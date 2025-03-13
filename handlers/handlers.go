package handlers

import (
	"dhdorr/story-point-poker/table"
	"fmt"
	"html/template"
	"net/http"
)

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("joining table: %v \n", r.Form)
	tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
	tmpl.Execute(w, nil)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("creating table: %v \n", r.Form)
	tml, err := table.GenerateTableSession(r.Form)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	fmt.Printf("session made: %v \n", tml)

	tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
	tmpl.Execute(w, nil)
}
