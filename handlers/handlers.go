package handlers

import (
	"net/http"
	"text/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home-wrapper.html"))
	tmpl.Execute(w, nil)
}

func HandleJoinTableMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/join-table-menu.html"))
	tmpl.Execute(w, nil)
}

func HandleCreateTableMenu(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create-table-menu.html"))
	tmpl.Execute(w, nil)
}
