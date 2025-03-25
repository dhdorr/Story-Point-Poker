package table_manager

import (
	"dhdorr/story-point-poker/table_session"
	"fmt"
	"html/template"
	"net/http"
)

type Table_Identifiers struct {
	TableID  string
	Passcode string
}

type Table_Map map[Table_Identifiers]table_session.Table_Session

type Table_Manager struct {
	TableMap Table_Map
}

func (tm *Table_Manager) HandleJoinTable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("recv join table form: %v", r.Form)
	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, nil)
}

func (tm *Table_Manager) HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("recv create table form: %v", r.Form)
	tmpl := template.Must(template.ParseFiles("templates/poker-table.html"))
	tmpl.Execute(w, nil)
}
