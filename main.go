package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Card struct {
	Value int
}

type Deck struct {
	Cards []Card
}

type Player struct {
	Username string
}

type Table_Template struct {
	Deck      []Card
	Skin      string
	TextColor string
	Players   []Player
}

type Session struct {
	ID          string
	Players     []Player
	CurrentDeck Deck
	Passcode    string
	Closed      bool
}

func NewSession(id, passcode string, players []Player, closed bool, deck Deck) *Session {
	return &Session{ID: id, Passcode: passcode, Players: players, Closed: closed, CurrentDeck: deck}
}

type Poker_Tables struct {
	active_sessions map[string]Session
}

func renderTemplate(w http.ResponseWriter, tmpl string, ts *Table_Template) {
	t, _ := template.ParseFiles("./templates/" + tmpl + ".html")
	t.Execute(w, ts)
}

func (poker *Poker_Tables) handle_test(w http.ResponseWriter, _ *http.Request) {
	data := Deck{
		Cards: []Card{
			{Value: 1},
			{Value: 2},
			{Value: 3},
			{Value: 4},
			{Value: 5},
		},
	}
	ts := Table_Template{Deck: data.Cards, Skin: "grey", TextColor: "black"}
	renderTemplate(w, "card", &ts)
}

func (poker *Poker_Tables) handle_join(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	session_id := r.URL.Query().Get("sessionID")
	passcode := r.URL.Query().Get("passcode")
	skin := r.Header.Get("bg_skin")
	text_color := r.Header.Get("bg_text")
	username := r.Header.Get("username")

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
	tmp.Players = append(tmp.Players, Player{Username: username})
	poker.active_sessions[session_id] = tmp

	fmt.Println(poker.active_sessions[session_id].Players)

	// renderTemplate(w, "test2", nil)
	ts := Table_Template{Deck: val.CurrentDeck.Cards, Skin: skin, TextColor: text_color, Players: tmp.Players}
	renderTemplate(w, "card", &ts)
}

func (poker *Poker_Tables) handle_create(w http.ResponseWriter, r *http.Request) {
	session_id := r.FormValue("sessionID")
	preset := r.FormValue("preset")
	passcode := r.FormValue("passcode")
	skin := r.Header.Get("bg_skin")
	text_color := r.Header.Get("bg_text")
	username := r.Header.Get("username")

	fmt.Printf("Session ID: %s | Preset: %s | Passcode: %s", session_id, preset, passcode)

	cards := make([]Card, 0, 8)

	if preset == "seq" {
		for i := 0; i < 8; i++ {
			cards = append(cards, Card{Value: i + 1})
		}
	} else if preset == "fib" {
		fib := []int{1, 2, 3, 5, 8, 13, 21, 34}
		for i := 0; i < 8; i++ {
			cards = append(cards, Card{Value: fib[i]})
		}
	}

	data := Deck{
		Cards: cards,
	}

	p := []Player{{Username: username}}

	poker.active_sessions[session_id] = *NewSession(session_id, passcode, p, false, data)

	tmp := poker.active_sessions[session_id]
	// tmp.Players = append(tmp.Players, Player{Username: username})
	// poker.active_sessions[session_id] = tmp

	fmt.Println(poker.active_sessions[session_id].Players)

	ts := Table_Template{Deck: data.Cards, Skin: skin, TextColor: text_color, Players: tmp.Players}
	renderTemplate(w, "card", &ts)
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
	poker_tables.active_sessions["abc123"] = Session{ID: "abc123", Players: []Player{Player{Username: "test"}}, CurrentDeck: data, Passcode: "test", Closed: false}
	// End Test

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index2.html")
	})

	http.HandleFunc("GET /test", poker_tables.handle_test)

	http.HandleFunc("GET /join", poker_tables.handle_join)

	http.HandleFunc("POST /create", poker_tables.handle_create)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
