package managers

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/table"
	"dhdorr/story-point-poker/templates"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

	ts.AddPlayerToTableSession(un)

	tm.Table_Sessions_M[t_id] = ts

	fmt.Printf("Session Joined: %v\n", ts.Players)

	// Check if max player count has been reached... REFACTOR ME
	tm.SetTableSessionState(t_id)

	tm.PrintTables()

	handlers.RenderTemplate(w, filename, ts)
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

	filename := "T-waiting.html"
	// w.Header().Add("tableID", ts.Table_ID)
	handlers.RenderTemplate(w, filename, *ts)
}

// called in the waiting for players screen to update client-side player count visuals
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

func (tm *Table_Manager) HandleStart(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	t_id := *table.NewTableSessionIdentifier(r.URL.Query().Get("tableID"), r.URL.Query().Get("passcode"))
	t := tm.Table_Sessions_M[t_id]
	t.Active_Round_ID, _ = strconv.Atoi(r.URL.Query().Get("activeRound"))
	tm.Table_Sessions_M[t_id] = t

	fmt.Printf("Table Started: %v\n", tm.Table_Sessions_M[t_id])

	filename := "T-poker-table.html"
	handlers.RenderTemplate(w, filename, tm.Table_Sessions_M[t_id])
}

func (tm *Table_Manager) HandleSelectCard(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	fmt.Printf("r.Formv: %v\n", r.FormValue("pcard"))
	fmt.Printf("r.Form: %v\n", r.Form)
	fmt.Println(r.Header.Get("tableID"))
	fmt.Println(r.Header.Get("passcode"))
	fmt.Println(r.Header.Get("username"))

	sv, _ := strconv.Atoi(r.FormValue("pcard"))
	t_id := table.Table_Session_Identifiers{Table_ID: r.Header.Get("tableID"), Passcode: r.Header.Get("passcode")}

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
	tmp.AddPlayerToTableSession(un)
	return tmp
}
