package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type PlayerBox struct {
	Name    string
	Value   int
	HasVal  bool
	IsAdmin bool
}

type Card struct {
	Value int
}

type Deck struct {
	Cards []Card
}

type Player struct {
	Username string
	isAdmin  bool
}

type Table_Template struct {
	Deck          []Card
	Skin          string
	TextColor     string
	Players       []Player
	TimeLimit     int
	TimeRemaining int
}

type Res_Card struct {
	Value       int
	VoteCount   int
	PlayerCount int
}

type Results_Template struct {
	Cards  []Res_Card
	Winner int
}

type Session struct {
	ID          string
	Players     []Player
	CurrentDeck Deck
	Passcode    string
	Closed      bool
	Choices     map[string]int
	TimeLimit   int
	StartTime   time.Time
}

func NewSession(id, passcode string, players []Player, closed bool, deck Deck, time_limit int) *Session {
	tmp := make(map[string]int)
	return &Session{ID: id, Passcode: passcode, Players: players, Closed: closed, CurrentDeck: deck, Choices: tmp, TimeLimit: time_limit, StartTime: time.Now()}
}

type Poker_Tables struct {
	active_sessions map[string]Session
}

func renderTemplate(w http.ResponseWriter, tmpl string, ts *Table_Template) {
	w.Header().Add("isAdmin", "true")
	t, _ := template.ParseFiles("./templates/" + tmpl + ".html")
	t.Execute(w, ts)
}

// func (poker *Poker_Tables) handle_test(w http.ResponseWriter, _ *http.Request) {
// 	data := Deck{
// 		Cards: []Card{
// 			{Value: 1},
// 			{Value: 2},
// 			{Value: 3},
// 			{Value: 4},
// 			{Value: 5},
// 		},
// 	}
// 	ts := Table_Template{Deck: data.Cards, Skin: "grey", TextColor: "black"}
// 	renderTemplate(w, "card", &ts)
// }

func (poker *Poker_Tables) handle_join(w http.ResponseWriter, r *http.Request) {
	session_id := r.FormValue("sessionID")
	passcode := r.FormValue("passcode")
	skin := r.Header.Get("bg_skin")
	text_color := r.Header.Get("bg_text")
	username := r.Header.Get("username")

	// fmt.Println(r.URL.Query())

	error_message := "SessionID or Passcode does not match"

	val, ok := poker.active_sessions[session_id]
	if !ok {
		fmt.Printf("No session with ID: %s", session_id)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, error_message)
		return
	}

	if val.Passcode != passcode {
		fmt.Printf("Passcode does not match for sessionID: %s", session_id)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, error_message)
		return
	}
	tmp := poker.active_sessions[session_id]
	tmp.Players = append(tmp.Players, Player{Username: username, isAdmin: false})
	poker.active_sessions[session_id] = tmp

	fmt.Println(poker.active_sessions[session_id].Players)

	// renderTemplate(w, "test2", nil)
	tr := time.Since(tmp.StartTime)
	fmt.Printf("tr: %v \n", tr)
	tr_1 := 30 - int(tr.Seconds())
	fmt.Printf("tr 1: %v \n\n", tr_1)

	ts := Table_Template{Deck: val.CurrentDeck.Cards, Skin: skin, TextColor: text_color, Players: tmp.Players, TimeLimit: tmp.TimeLimit, TimeRemaining: tr_1}
	renderTemplate(w, "card", &ts)
}

func (poker *Poker_Tables) handle_results(w http.ResponseWriter, r *http.Request) {
	session_id := r.Header.Get("sessionID")

	tmp := poker.active_sessions[session_id]

	res_choices := make(map[int]int)
	for _, b := range tmp.Choices {
		check, ok := res_choices[b]
		if !ok {
			res_choices[b] = 1
		} else {
			check += 1
			res_choices[b] = check
		}
	}

	prev_high := 0
	high_idx := 0
	for k, v := range res_choices {
		if v > prev_high {
			prev_high = v
			high_idx = k
		}
	}

	data := Results_Template{
		Cards:  []Res_Card{},
		Winner: high_idx,
	}
	for _, c := range tmp.CurrentDeck.Cards {
		new_res_card := Res_Card{Value: c.Value, VoteCount: res_choices[c.Value], PlayerCount: len(tmp.Players)}
		data.Cards = append(data.Cards, new_res_card)
	}
	fmt.Println(data)

	// data := Results_Template{
	// 	Cards: []Res_Card{
	// 		{Value: 1},
	// 		{Value: 2},
	// 		{Value: 3},
	// 		{Value: 4},
	// 		{Value: 5},
	// 	},
	// }
	t, _ := template.ParseFiles("./templates/results.html")
	t.Execute(w, data)
}

func (poker *Poker_Tables) handle_create(w http.ResponseWriter, r *http.Request) {
	session_id := r.FormValue("sessionID")
	preset := r.FormValue("preset")
	passcode := r.FormValue("passcode")
	num_cards := r.FormValue("num-select")
	time_limit := r.FormValue("time-select")
	skin := r.Header.Get("bg_skin")
	text_color := r.Header.Get("bg_text")
	username := r.Header.Get("username")

	fmt.Printf("Session ID: %s | Preset: %s | Passcode: %s | Number of Cards: %s | Time Limit: %s \n", session_id, preset, passcode, num_cards, time_limit)

	num, err := strconv.Atoi(num_cards)
	if err != nil {
		num = 8
		fmt.Printf("There was an error parsing the num-select list: %s\n", err)
	}

	if num > 12 {
		num = 12
	} else if num < 6 {
		num = 6
	}

	tl, err := strconv.Atoi(time_limit)
	if err != nil {
		tl = 30
		fmt.Printf("There was an error parsing the time limit list: %s\n", err)
	}

	if tl > 90 {
		tl = 90
	} else if tl < 10 {
		tl = 10
	}

	cards := make([]Card, 0, num)

	if preset == "seq" {
		for i := 0; i < num; i++ {
			cards = append(cards, Card{Value: i + 1})
		}
	} else if preset == "fib" {
		prev := 1
		val := 1
		for i := 0; i < num; i++ {
			fmt.Printf("fib - prev: [%v], val [%v]\n", prev, val)
			cards = append(cards, Card{Value: val})

			tmp := val + prev
			prev = val
			val = tmp

		}
	}

	tl_val, err := strconv.Atoi(time_limit)
	if err != nil {
		fmt.Println(err)
		tl_val = 15
	}
	// infinity option
	if tl_val == -1 {
		tl_val = 600
	}

	fmt.Printf("Time Limit: %v \n", tl_val)

	timer1 := time.NewTimer(time.Duration(tl_val) * time.Second)
	go func() {
		<-timer1.C
		fmt.Println("Session Over: Time Expired...")
		tmp := poker.active_sessions[session_id]
		tmp.Closed = true
		poker.active_sessions[session_id] = tmp
	}()

	data := Deck{
		Cards: cards,
	}

	p := []Player{{Username: username, isAdmin: true}}

	poker.active_sessions[session_id] = *NewSession(session_id, passcode, p, false, data, tl_val)

	tmp := poker.active_sessions[session_id]

	fmt.Println(poker.active_sessions[session_id].Players)

	ts := Table_Template{Deck: data.Cards, Skin: skin, TextColor: text_color, Players: tmp.Players, TimeLimit: tmp.TimeLimit, TimeRemaining: tmp.TimeLimit}
	renderTemplate(w, "card-admin", &ts)
}

func (poker *Poker_Tables) handle_choose(w http.ResponseWriter, r *http.Request) {
	session_id := r.Header.Get("sessionID")
	username := r.Header.Get("username")
	value := r.FormValue("value")
	v, _ := strconv.Atoi(value)
	fmt.Println("value: ", v)
	fmt.Println("sessionID: ", session_id)

	tmp := poker.active_sessions[session_id]
	if tmp.Closed {
		fmt.Println("Selection Cancelled, round is over")
		return
	}
	tmp.Choices[username] = v
	poker.active_sessions[session_id] = tmp

	fmt.Println(poker.active_sessions[session_id].Choices)

	t, _ := template.ParseFiles("./templates/card-selected.html")
	type cardSelected struct {
		Value int
	}
	data := cardSelected{Value: v}
	t.Execute(w, data)

}

func (poker *Poker_Tables) handle_player_box(w http.ResponseWriter, r *http.Request) {
	session_id := r.Header.Get("sessionID")
	inc_username := r.Header.Get("username")

	tmp := poker.active_sessions[session_id]
	pb_arr := make([]PlayerBox, 0)

	show_vals := false

	for _, v := range tmp.Players {
		cv, ok := tmp.Choices[v.Username]
		if !ok {
			cv = -1
		}
		// shows the admin everyone's selection
		if v.Username == inc_username {
			show_vals = v.isAdmin
		}
		pb := PlayerBox{Name: v.Username, Value: cv, HasVal: ok, IsAdmin: show_vals}
		pb_arr = append(pb_arr, pb)
	}

	t, _ := template.ParseFiles("./templates/player-box.html")

	t.Execute(w, pb_arr)
}

func main() {
	poker_tables := &Poker_Tables{active_sessions: make(map[string]Session)}

	// Test Values
	data := Deck{
		Cards: []Card{
			{Value: 1},
			{Value: 2},
			{Value: 3},
			{Value: 4},
			{Value: 5},
		},
	}
	tmp := make(map[string]int)
	poker_tables.active_sessions["abc123"] = Session{ID: "abc123", Players: []Player{{Username: "test", isAdmin: true}}, CurrentDeck: data, Passcode: "test", Closed: false, Choices: tmp, TimeLimit: 60}
	// End Test

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index2.html")
	})

	// http.HandleFunc("GET /test", poker_tables.handle_test)

	http.HandleFunc("GET /results", poker_tables.handle_results)

	http.HandleFunc("POST /join", poker_tables.handle_join)

	http.HandleFunc("POST /create", poker_tables.handle_create)

	http.HandleFunc("POST /choose", poker_tables.handle_choose)

	http.HandleFunc("GET /playerBox", poker_tables.handle_player_box)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
