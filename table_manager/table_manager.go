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

	data := templates.Waiting{MaxPlayers: t_map[t_id].Settings.MaxPlayers, PlayerCount: len(t_map[t_id].Players)}
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

	data := templates.Waiting{MaxPlayers: mp, PlayerCount: len(pl)}
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

func (tm *Table_Map) HandlePlayerCount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	p_count := len(t_map[t_id].Players)
	tmpl, _ := template.New("count").Parse(strconv.Itoa(p_count))
	tmpl.Execute(w, nil)
}

func (tm *Table_Map) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	// un := r.FormValue("username")
	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")

	t_id := Table_Identifiers{TableID: id, Passcode: pc}
	t_map := *tm
	t_cards := t_map[t_id].Cards

	data := templates.Game_Table{Cards: t_cards}
	tmpl, _ := template.ParseFiles("templates/game-table.html")
	tmpl.Execute(w, data)
}

// func (tm *Table_Map) ChangeActiveRound(t_id Table_Identifiers, roundID int) {
// 	t_map := *tm
// 	table := t_map[t_id]

// 	table.ActiveRoundID = roundID

// 	t_map[t_id] = table
// }
