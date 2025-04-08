package table_manager

import (
	"dhdorr/story-point-poker/card"
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/round"
	"dhdorr/story-point-poker/table_session"
	"dhdorr/story-point-poker/templates"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Table_Identifiers struct {
	TableID  string
	Passcode string
}

type Table_Map map[Table_Identifiers]table_session.Table

func (tm *Table_Map) HandleJoinTable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")
	t_id := Table_Identifiers{
		TableID:  id,
		Passcode: pc,
	}

	un := r.FormValue("username")

	t_map := *tm
	_, ok := t_map[t_id]
	if !ok {
		w.WriteHeader(404)
		fmt.Fprintf(w, "table does not exist! id: %v, pc: %v", id, pc)
		return
	}

	t_state := t_map[t_id].State
	if t_state == table_session.StateClosed {
		fmt.Fprintf(w, "table is at capacity! id: %v, pc: %v", id, pc)
		return
	}

	table := t_map[t_id]
	table.Players = append(table.Players, player.Player{PlayerID: "test-guest", Username: un})
	if len(table.Players) >= table.Settings.MaxPlayers {
		table.State = table_session.StateClosed
	}
	t_map[t_id] = table

	d_style := template.CSS("display: none;")
	data := templates.Waiting{
		MaxPlayers:  t_map[t_id].Settings.MaxPlayers,
		PlayerCount: len(t_map[t_id].Players),
		TableID:     id,
		Passcode:    pc,
		Username:    un,
		Ready:       false,
		Style:       d_style,
	}
	fmt.Printf("data: %v\n", data)
	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, data)
}

func (tm *Table_Map) HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")
	t_id := Table_Identifiers{
		TableID:  id,
		Passcode: pc,
	}

	t_map := *tm
	_, ok := t_map[t_id]
	if ok {
		fmt.Fprintf(w, "table already exists, cannot overwrite! id: %v, pc: %v", id, pc)
		return
	}

	un := r.FormValue("username")
	cl := r.FormValue("cardLayout")
	nc, err := strconv.Atoi(r.FormValue("numCards")) // int
	if err != nil {
		fmt.Fprintf(w, "invalid value for numCards: %v\n", nc)
		return
	}
	nr, err := strconv.Atoi(r.FormValue("numRounds")) // int
	if err != nil {
		fmt.Fprintf(w, "invalid value for numRounds: %v\n", nr)
		return
	}
	mp, err := strconv.Atoi(r.FormValue("maxPlayers")) // int
	if err != nil {
		fmt.Fprintf(w, "invalid value for maxPlayers: %v\n", mp)
		return
	}
	rtl, err := strconv.Atoi(r.FormValue("roundTimeLimit")) // int
	if err != nil {
		fmt.Fprintf(w, "invalid value for roundTimeLimit: %v\n", rtl)
		return
	}
	ttl := 1800                   // 30 minutes in seconds
	ar := 0                       // first round
	st := table_session.StateOpen // initial table state

	ad := player.Player{
		PlayerID: "test-admin",
		Username: un,
	}
	pl := make([]player.Player, 0, mp) // initial empty slice of players
	pl = append(pl, ad)

	rd := make([]round.Round, 0, nr) // initial empty slice of rounds
	for i := 0; i < nr; i++ {
		rd = append(rd, round.Round{Phase: round.PhaseWaiting, Votes: make([]card.Vote, 0, mp)})
	}

	cd := make([]card.Card, 0, nc) // initial empty slice of cards
	val := 1
	prev := 1
	for i := range nc {
		if cl == "seq" {
			val = i + 1
			cd = append(cd, card.Card{Value: val})
		} else if cl == "fib" {
			cd = append(cd, card.Card{Value: val})
			val = val + prev
			prev = val - prev
		}
	}

	stime := time.Now()
	// eTime := time.Time{} // zero value for end time

	settings := table_session.Settings{
		CardLayout:     cl,
		NumCards:       nc,
		NumRounds:      nr,
		MaxPlayers:     mp,
		RoundTimeLimit: rtl,
		TableTimeLimit: ttl,
	}

	new_table := table_session.Table{
		TableID:       id,
		Passcode:      pc,
		ActiveRoundID: ar,
		State:         st,
		Players:       pl,
		Rounds:        rd,
		Cards:         cd,
		Admin:         ad,
		StartTime:     stime,
		Settings:      settings,
	}

	t_map[t_id] = new_table

	data := templates.Waiting{
		MaxPlayers:  mp,
		PlayerCount: len(pl),
		TableID:     id,
		Passcode:    pc,
		Username:    un,
		Ready:       false,
	}
	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, data)
}

func (tm *Table_Map) HandleDeleteTables(w http.ResponseWriter, r *http.Request) {
	t_map := *tm

	for k := range t_map {
		delete(t_map, k)
	}

	fmt.Fprintf(w, "Table Map has been cleared!")
}

func (tm *Table_Map) HandleWaitingUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	un := r.URL.Query().Get("username")

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	p_count := len(t_map[t_id].Players)
	p_max := t_map[t_id].Settings.MaxPlayers
	ready := false

	fmt.Println(t_map[t_id].Rounds[t_map[t_id].ActiveRoundID].Phase)
	if t_map[t_id].Rounds[t_map[t_id].ActiveRoundID].Phase == round.PhaseStarted {
		fmt.Println("round is set to started...")
		ready = true
	}

	data := templates.Waiting{
		MaxPlayers:  p_max,
		PlayerCount: p_count,
		TableID:     id,
		Passcode:    pc,
		Username:    un,
		Ready:       ready,
	}

	tmpl, _ := template.ParseFiles("templates/waiting-update.html")
	tmpl.Execute(w, data)
}

func (tm *Table_Map) HandleGameUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm

	table := t_map[t_id]
	t_rounds := table.Rounds
	t_r := t_rounds[table.ActiveRoundID]

	is_done := false
	if t_r.Phase == round.PhaseFinished {
		is_done = true
	}

	data := templates.Game_Update{
		IsDone: is_done,
	}

	tmpl, _ := template.ParseFiles("templates/game-update.html")
	tmpl.Execute(w, data)
}

func (tm *Table_Map) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	un := r.FormValue("username")
	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")

	fmt.Printf("id: %v, pc: %v, un: %v\n", id, pc, un)

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	t_cards := t_map[t_id].Cards

	if t_map[t_id].Admin.Username == un {
		table := t_map[t_id]
		t_rounds := table.Rounds
		t_r := t_rounds[table.ActiveRoundID]
		t_r.Phase = round.PhaseStarted // set round phase to Started
		t_rounds[table.ActiveRoundID] = t_r
		table.Rounds = t_rounds
		t_map[t_id] = table

		timer1 := time.NewTimer(30 * time.Second)
		go func(tm1 *Table_Map) {
			<-timer1.C

			t_id1 := Table_Identifiers{TableID: id, Passcode: pc}
			t_map1 := *tm1
			table1 := t_map1[t_id1]
			t_rounds1 := table1.Rounds
			t_r1 := t_rounds1[table1.ActiveRoundID]
			t_r1.Phase = round.PhaseFinished // set round phase to Ended
			t_rounds1[table1.ActiveRoundID] = t_r1
			table1.Rounds = t_rounds1
			t_map1[t_id1] = table1
		}(tm)
	}

	fmt.Printf("t_map: %v\n", t_map)

	data := templates.Game_Table{
		Cards:    t_cards,
		TableID:  id,
		Passcode: pc,
		Username: un,
		IsDone:   false,
	}

	tmpl, _ := template.ParseFiles("templates/game-table.html")
	tmpl.Execute(w, data)
}

func (tm *Table_Map) HandleVote(w http.ResponseWriter, r *http.Request) {
	un := r.FormValue("username")
	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")
	val, _ := strconv.Atoi(r.FormValue("cardValue"))

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	table := t_map[t_id]

	cards := table.Cards
	vote := card.Vote{PlayerID: un}
	for _, c := range cards {
		if c.Value == val {
			vote.Card = c
		}
	}

	ar := table.ActiveRoundID
	rounds := table.Rounds
	rd := rounds[ar]

	is_new_vote := true
	for i, v := range rd.Votes {
		if vote.PlayerID == v.PlayerID {
			rd.Votes[i] = vote
			is_new_vote = false
			fmt.Printf("Player has already voted: %v\n", un)
			break
		}
	}
	if is_new_vote {
		rd.Votes = append(rd.Votes, vote)
	}

	rounds[ar] = rd
	table.Rounds = rounds
	t_map[t_id] = table

	fmt.Printf("t_map[t_id]: %v\n", t_map[t_id])

	tmpl, _ := template.New("count").Parse("picked")
	tmpl.Execute(w, nil)
}

func (tm *Table_Map) HandleRoundResults(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	table := t_map[t_id]

	t_rounds := table.Rounds
	t_r := t_rounds[table.ActiveRoundID]

	r_cards := make([]templates.Results_Card, 0, len(table.Cards))
	for _, v := range table.Cards {
		temp := templates.Results_Card{
			Value: v.Value,
		}
		r_cards = append(r_cards, temp)
	}

	for j, v := range r_cards {
		for _, vt := range t_r.Votes {
			if vt.Card.Value == v.Value {
				r_cards[j].Votes += 1
			}
		}
	}

	data := templates.Round_Results{
		Cards: r_cards,
	}

	tmpl, _ := template.ParseFiles("templates/results.html")
	tmpl.Execute(w, data)
}

// func (tm *Table_Map) ChangeActiveRound(t_id Table_Identifiers, roundID int) {
// 	t_map := *tm
// 	table := t_map[t_id]

// 	table.ActiveRoundID = roundID

// 	t_map[t_id] = table
// }
