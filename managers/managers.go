package managers

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/table"
	"dhdorr/story-point-poker/templates"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Table_Manager struct {
	Table_Sessions_M table.Table_Map
}

func (tm *Table_Manager) PrintTables() {
	fmt.Printf("tm.Table_Sessions_M: %v\n", tm.Table_Sessions_M)
}

func (tm *Table_Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	t_id := table.Table_Session_Identifiers{Table_ID: r.FormValue("tableID"), Passcode: r.FormValue("passcode")}
	un := r.FormValue("username")

	ts, ok := tm.Table_Sessions_M[t_id]
	if !ok {
		fmt.Fprintf(w, "Unable to find table: %v", t_id.Table_ID)
		return
	}

	filename := "T-poker-table.html"
	// TODO:
	if ts.Session_State == table.StateClosed {
		fmt.Fprintf(w, "Unable to join, table closed: %v", t_id.Table_ID)
		return
	}

	if ts.Rounds[ts.Active_Round_ID].Phase == table.PhaseWaitingForPlayers {
		filename = "T-waiting.html"
	}
	// ^^^

	ts.AddPlayerToTableSession(un, false)

	tm.Table_Sessions_M[t_id] = ts

	fmt.Printf("Session Joined: %v\n", ts.Players)

	// Check if max player count has been reached... REFACTOR ME
	tm.SetTableSessionState(t_id)

	tm.PrintTables()

	data := templates.Waiting_For_Players{
		Table_ID:    ts.Table_ID,
		Passcode:    ts.Passcode,
		PlayerCount: len(ts.Players),
		MaxPlayers:  ts.Settings.Player_Max,
		IsAdmin:     false,
		IsReady:     0,
		Username:    r.FormValue("username"),
		TimeLimit:   ts.Settings.Round_Time_Limit,
	}

	handlers.RenderTemplate(w, filename, data)
}

func (tm *Table_Manager) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ts, err := handlers.HandleCreate(w, r)
	if err != nil {
		fmt.Fprintf(w, "Failed to create a new table session: %v", err)
		return
	}
	t_id := table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}

	// tm.AddNewTableSession(t_id, ts)
	tm.Table_Sessions_M[t_id] = *ts

	// Check if max player count has been reached... REFACTOR ME
	tm.SetTableSessionState(t_id)

	tm.PrintTables()

	data := templates.Waiting_For_Players{
		Table_ID:    ts.Table_ID,
		Passcode:    ts.Passcode,
		PlayerCount: len(ts.Players),
		MaxPlayers:  ts.Settings.Player_Max,
		IsAdmin:     true,
		IsReady:     0,
		Username:    r.FormValue("username"),
		TimeLimit:   ts.Settings.Round_Time_Limit,
	}

	filename := "T-waiting.html"
	// w.Header().Add("tableID", ts.Table_ID)
	handlers.RenderTemplate(w, filename, data)
}

// called in the waiting for players screen to update client-side player count visuals
func (tm *Table_Manager) HandlePlayerCount(w http.ResponseWriter, r *http.Request) {
	t_id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	ts := tm.Table_Sessions_M[*table.NewTableSessionIdentifier(t_id, pc)]
	ts.PrintTableSessionPlayers()

	isready := 0
	fmt.Println(ts.Rounds[ts.Active_Round_ID].Phase)
	if ts.Rounds[ts.Active_Round_ID].Phase == table.PhaseStarted {
		fmt.Println("round ready")
		isready = 1
	}

	data := templates.Player_Count{PlayerCount: len(ts.Players), IsReady: isready}
	filename := "T-player-count.html"
	handlers.RenderTemplate(w, filename, data)
}

func (tm *Table_Manager) HandleCheckForNewPlayers(w http.ResponseWriter, r *http.Request) {
	t_id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	ts := tm.Table_Sessions_M[*table.NewTableSessionIdentifier(t_id, pc)]
	ts.PrintTableSessionPlayers()

	// if player max is reached...

	filename := "T-players.html"
	handlers.RenderTemplate(w, filename, ts)
}

func (tm *Table_Manager) HandleCheckForNewVotes(w http.ResponseWriter, r *http.Request) {
	t_id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	ts := tm.Table_Sessions_M[*table.NewTableSessionIdentifier(t_id, pc)]

	rd := ts.Rounds[0]
	filename := "T-votes.html"
	handlers.RenderTemplate(w, filename, rd)
}

func (tm *Table_Manager) HandleStart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	t_id := *table.NewTableSessionIdentifier(r.URL.Query().Get("tableID"), r.URL.Query().Get("passcode"))
	t := tm.Table_Sessions_M[t_id]

	player := player.Player{}
	for _, p := range t.Players {
		if p.Username == r.URL.Query().Get("username") {
			player = p
			break
		}
	}

	if player.IsAdmin {
		// t.Active_Round_ID, _ = strconv.Atoi(r.URL.Query().Get("activeRound"))
		t.Rounds[t.Active_Round_ID].Phase = table.PhaseStarted
		fmt.Printf("round phase transition: %v\n", t.Rounds[t.Active_Round_ID].Phase)
		t.Rounds[t.Active_Round_ID].Start_Time = time.Now()

		round_timer := time.NewTimer(time.Duration(t.Settings.Round_Time_Limit) * time.Second)
		go func(t_id2 table.Table_Session_Identifiers) {
			<-round_timer.C
			temp := tm.Table_Sessions_M[t_id2]
			temp.Rounds[temp.Active_Round_ID].Phase = table.PhaseFinished
			fmt.Println("round is over")
			tm.Table_Sessions_M[t_id2] = temp
			fmt.Printf("tm.Table_Sessions_M[t_id2].Rounds[tm.Table_Sessions_M[t_id2].Active_Round_ID].Phase: %v\n", tm.Table_Sessions_M[t_id2].Rounds[tm.Table_Sessions_M[t_id2].Active_Round_ID].Phase)
		}(t_id)
	}

	tm.Table_Sessions_M[t_id] = t

	fmt.Printf("Table Started: %v\n", tm.Table_Sessions_M[t_id])

	filename := "T-poker-table.html"

	data := templates.Game_Table{
		Cards:    t.Cards,
		Players:  t.Players,
		Table_ID: t.Table_ID,
		Passcode: t.Passcode,
		IsAdmin:  player.IsAdmin,
		Username: player.Username,
	}
	handlers.RenderTemplate(w, filename, data)
}

func (tm *Table_Manager) HandleEndRound(w http.ResponseWriter, r *http.Request) {
	t_id := *table.NewTableSessionIdentifier(r.URL.Query().Get("tableID"), r.URL.Query().Get("passcode"))
	ts := tm.Table_Sessions_M[t_id]

	un := r.URL.Query().Get("username")

	isAdmin, _ := strconv.ParseBool(r.URL.Query().Get("isAdmin"))

	rc_arr := make([]table.Results_Card, 0, len(ts.Cards))
	for _, v := range ts.Cards {
		rc_arr = append(rc_arr, table.Results_Card{Value: v.Value, VoteCount: 0, Winner: false})
	}
	for _, v2 := range ts.Rounds[ts.Active_Round_ID].Votes {
		for i := range rc_arr {
			if rc_arr[i].Value == v2 {
				rc_arr[i].VoteCount += 1
			}
		}
	}
	most_voted_index := 0
	for i := 1; i < len(rc_arr); i++ {
		if rc_arr[i].VoteCount > rc_arr[i-1].VoteCount {
			most_voted_index = i
		}
	}
	rc_arr[most_voted_index].Winner = true

	data := templates.Results{
		Cards:       rc_arr,
		ActiveRound: ts.Active_Round_ID,
		NextRound:   ts.Active_Round_ID + 1,
		IsAdmin:     isAdmin,
		Table_ID:    t_id.Table_ID,
		Passcode:    t_id.Passcode,
		Username:    un,
	}

	handlers.RenderTemplate(w, "T-results.html", data)
}

func (tm *Table_Manager) HandleNextRound(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t_id := table.Table_Session_Identifiers{Table_ID: r.FormValue("tableID"), Passcode: r.FormValue("passcode")}
	un := r.FormValue("username")

	isAdmin, _ := strconv.ParseBool(r.FormValue("isAdmin"))

	ts := tm.Table_Sessions_M[t_id]

	if ts.Active_Round_ID+1 >= ts.Settings.Number_Of_Rounds {
		fmt.Fprintf(w, "game over, no more rounds!")
		return
	}

	if isAdmin {

		ts.Active_Round_ID = ts.Active_Round_ID + 1
		ts.Rounds[ts.Active_Round_ID].Phase = table.PhaseStarted
		tm.Table_Sessions_M[t_id] = ts
	}

	data := templates.Game_Table{
		Cards:    ts.Cards,
		Players:  ts.Players,
		Table_ID: ts.Table_ID,
		Passcode: ts.Passcode,
		IsAdmin:  isAdmin,
		Username: un,
	}

	handlers.RenderTemplate(w, "T-poker-table.html", data)
}

func (tm *Table_Manager) HandleCheckForRoundChange(w http.ResponseWriter, r *http.Request) {
	t_id := *table.NewTableSessionIdentifier(r.URL.Query().Get("tableID"), r.URL.Query().Get("passcode"))
	ts := tm.Table_Sessions_M[t_id]

	r_id, _ := strconv.Atoi(r.URL.Query().Get("activeRound"))
	val := false

	fmt.Printf("inc_roundID: %v, ts_roundID: %v\n", r_id, ts.Active_Round_ID)
	if ts.Active_Round_ID != r_id {
		val = true
	}

	data := templates.Should_Change_Round{
		ChangeRound: val,
	}

	handlers.RenderTemplate(w, "T-change-round.html", data)
}

func (tm *Table_Manager) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("r.Formv: %v\n", r.FormValue("pcard"))
	// fmt.Printf("r.Form: %v\n", r.Form)
	// fmt.Println(r.Header.Get("tableID"))
	// fmt.Println(r.Header.Get("passcode"))
	// fmt.Println(r.Header.Get("username"))

	sv, _ := strconv.Atoi(r.FormValue("pcard"))
	t_id := table.Table_Session_Identifiers{Table_ID: r.Header.Get("tableID"), Passcode: r.Header.Get("passcode")}

	// if tm.Table_Sessions_M[t_id].Rounds[tm.Table_Sessions_M[t_id].Active_Round_ID].Phase == table.PhaseFinished {
	// 	fmt.Fprintf(w, "Round is over: %v", sv)
	// 	return
	// }

	valid_choice := false
	for _, v := range tm.Table_Sessions_M[t_id].Cards {
		if sv == v.Value {
			valid_choice = true
		}
	}

	if !valid_choice {
		fmt.Fprintf(w, "Bad choice: %v", sv)
		return
	}

	ar_id := tm.Table_Sessions_M[t_id].Active_Round_ID
	for i, v := range tm.Table_Sessions_M[t_id].Players {
		if v.Username == r.Header.Get("username") {
			tm.Table_Sessions_M[t_id].Rounds[ar_id].SubmitVote(&tm.Table_Sessions_M[t_id].Players[i], sv)
		}
	}

	// tm.Table_Sessions_M[t_id].Rounds[0].SubmitVote(&tm.Table_Sessions_M[t_id].Players[0], sv)

	fmt.Println(tm.Table_Sessions_M[t_id].Rounds[0])
	data := templates.Gen_Test_A{Value: sv}
	filename := "T-card.html"
	handlers.RenderTemplate(w, filename, data)
}

func (tm *Table_Manager) AddNewTableSession(t_id table.Table_Session_Identifiers, ts *table.Table_Session) {
	tm.Table_Sessions_M[t_id] = *ts
}

func (tm *Table_Manager) SetTableSessionState(t_id table.Table_Session_Identifiers) {
	ts := tm.Table_Sessions_M[t_id]
	new_state := table.StateOpen
	if len(ts.Players) >= ts.Settings.Player_Max {
		new_state = table.StateClosed
	}

	ts_new := table.Table_Session{
		Table_ID:        ts.Table_ID,
		Passcode:        ts.Passcode,
		Settings:        ts.Settings,
		Players:         ts.Players,
		Rounds:          ts.Rounds,
		Active_Round_ID: ts.Active_Round_ID,
		Session_State:   new_state,
		Cards:           ts.Cards,
	}
	tm.Table_Sessions_M[t_id] = ts_new
}

// TESTING

func (tm *Table_Manager) HandleTest() {
	test := url.Values{
		"tableID":        []string{"test_id"},
		"passcode":       []string{"test_pc"},
		"cardLayout":     []string{"seq"},
		"username":       []string{"test_un"},
		"numCards":       []string{"6"},
		"numRounds":      []string{"1"},
		"roundTimeLimit": []string{"30"},
		"playerMax":      []string{"10"},
	}
	ts, err := HandleTestCreate(test)
	if err != nil {
		fmt.Println("Failed to create a test table session")
		return
	}

	tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}] = *ts

	ts_j := tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}]

	tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}] = *HandleTestJoin(test, &ts_j)
}

func HandleTestCreate(form_values url.Values) (*table.Table_Session, error) {
	ts, err := handlers.GenerateTableSession(form_values)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func HandleTestJoin(form_values url.Values, tmp *table.Table_Session) *table.Table_Session {
	un := form_values.Get("username")
	tmp.AddPlayerToTableSession(un, true)
	return tmp
}
