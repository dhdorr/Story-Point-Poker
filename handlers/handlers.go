package handlers

import (
	"dhdorr/story-point-poker/table"
	"dhdorr/story-point-poker/templates"
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

func HandleCreate(w http.ResponseWriter, r *http.Request) (*table.Table_Session, error) {
	r.ParseForm()
	fmt.Printf("creating table: %v \n", r.Form)
	ts, err := table.GenerateTableSession(r.Form)
	if err != nil {
		return nil, err
	}

	fmt.Printf("session made: %v \n", ts)
	return ts, nil
}

func RenderTemplate[V templates.Gen_Table_Session_Interface](w http.ResponseWriter, data V) {
	tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
	tmpl.Execute(w, data)
}
