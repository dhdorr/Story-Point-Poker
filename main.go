package main

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	return &Page{Title: title, Body: nil}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("./templates/" + tmpl + ".html")
	t.Execute(w, p)
}

func main() {

	// tmpl := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		title := "Derek?"
		p, _ := loadPage(title)
		renderTemplate(w, "test", p)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", nil)
}
