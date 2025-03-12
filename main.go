package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to Story Point Poker 2")

	http.Handle("/", http.FileServer(http.Dir(".")))

	fs_pages := http.FileServer(http.Dir("pages"))
	http.HandleFunc("GET /pages/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/pages/", fs_pages).ServeHTTP(w, r)
	})

	// serve dynamic html templates
	fs_templates := http.FileServer(http.Dir("templates"))
	http.HandleFunc("GET /templates/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
		http.StripPrefix("/templates/", fs_templates).ServeHTTP(w, r)
	})

	http.HandleFunc("POST /joinTable", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Printf("joining table: %v \n", r.Form)
		tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
		tmpl.Execute(w, nil)
	})

	// serve css and js, from html pages
	fs_static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
