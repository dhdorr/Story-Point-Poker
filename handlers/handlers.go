package handlers

import (
	"dhdorr/story-point-poker/table"
	"dhdorr/story-point-poker/templates"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

func HandleCreate(w http.ResponseWriter, r *http.Request) (*table.Table_Session, error) {
	r.ParseForm()
	fmt.Printf("creating table: %v \n", r.Form)
	ts, err := GenerateTableSession(r.Form)
	if err != nil {
		return nil, err
	}

	// add player to new session, after it has been created
	ts.AddPlayerToTableSession(r.FormValue("username"))

	fmt.Printf("session made: %v \n", ts)
	return ts, nil
}

func RenderTemplate[V templates.Gen_Table_Session_Interface](w http.ResponseWriter, filename string, data V) {
	tmpl := template.Must(template.ParseFiles("templates/" + filename))
	tmpl.Execute(w, data)
}

func GenerateTableSession(form_values url.Values) (*table.Table_Session, error) {
	tsc := table.Table_Session_Constructor{}
	tsc.ID = form_values.Get("tableID")
	tsc.PC = form_values.Get("passcode")
	tsc.CL = form_values.Get("cardLayout")
	tsc.UN = form_values.Get("username")

	nc, err := strconv.Atoi(form_values.Get("numCards"))
	if err != nil {
		return nil, err
	}
	tsc.NC = nc

	nr, err := strconv.Atoi(form_values.Get("numRounds"))
	if err != nil {
		return nil, err
	}
	tsc.NR = nr

	tl, err := strconv.Atoi(form_values.Get("roundTimeLimit"))
	if err != nil {
		return nil, err
	}
	tsc.TL = tl

	pm, err := strconv.Atoi(form_values.Get("playerMax"))
	if err != nil {
		return nil, err
	}
	tsc.PM = pm

	tsc.AR = -1

	return table.NewTableSessionConstructed(tsc), nil
}
