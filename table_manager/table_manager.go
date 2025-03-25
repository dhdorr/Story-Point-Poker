package table_manager

import (
	"dhdorr/story-point-poker/card"
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/round"
	"dhdorr/story-point-poker/table_session"
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
	fmt.Printf("recv join table form: %v\n", r.Form)

	id := r.FormValue("tableID")
	pc := r.FormValue("passcode")
	t_id := Table_Identifiers{
		TableID:  id,
		Passcode: pc,
	}

	t_map := *tm
	_, ok := t_map[t_id]
	if !ok {
		w.WriteHeader(404)
		fmt.Fprintf(w, "table does not exist! id: %v, pc: %v", id, pc)
		return
	}

	fmt.Printf("tm: %v\n", tm)

	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, nil)
}

func (tm *Table_Map) HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("recv create table form: %v\n", r.Form)

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
	fmt.Printf("cd: %v\n", cd)

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

	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, nil)
}
