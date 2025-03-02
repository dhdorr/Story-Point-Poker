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

type Session struct {
	ID          string
	Players     int
	CurrentDeck Deck
}

func renderTemplate(w http.ResponseWriter, tmpl string, d *Deck) {
	t, _ := template.ParseFiles("./templates/" + tmpl + ".html")
	t.Execute(w, d)
}

func main() {
	active_sessions := make(map[string]Session)

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
	active_sessions["abc123"] = Session{ID: "abc123", Players: 0, CurrentDeck: data}
	// End Test

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index2.html")
	})

	http.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		data := Deck{
			Cards: []Card{
				{Value: 1},
				{Value: 2},
				{Value: 3},
				{Value: 4},
				{Value: 5},
			},
		}
		renderTemplate(w, "card", &data)
	})

	http.HandleFunc("GET /join", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Query())
		session_id := r.URL.Query().Get("sessionID")
		val, ok := active_sessions[session_id]
		if !ok {
			fmt.Printf("No session with ID: %s", session_id)
			w.WriteHeader(404)
			fmt.Fprintf(w, "No session found with ID: %s", session_id)
			return
		}
		// renderTemplate(w, "test2", nil)
		renderTemplate(w, "card", &val.CurrentDeck)
	})

	http.HandleFunc("POST /create", func(w http.ResponseWriter, r *http.Request) {
		session_id := r.FormValue("sessionID")
		preset := r.FormValue("preset")
		fmt.Printf("Session ID: %s | Preset: %s", session_id, preset)

		cards := []Card{}

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
		renderTemplate(w, "card", &data)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
