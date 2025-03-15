package handlers

import (
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/table"
	"dhdorr/story-point-poker/templates"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

func HandleJoin(w http.ResponseWriter, r *http.Request, tsm *table.Table_Map) (*table.Table_Session, error) {
	r.ParseForm()
	fmt.Printf("joining table: %v \n", r.Form)

	a, err := AddPlayerToTableSession(r.Form, tsm)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func HandleCreate(w http.ResponseWriter, r *http.Request) (*table.Table_Session, error) {
	r.ParseForm()
	fmt.Printf("creating table: %v \n", r.Form)
	ts, err := GenerateTableSession(r.Form)
	if err != nil {
		return nil, err
	}

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

func AddPlayerToTableSession(form_values url.Values, table_sessions *table.Table_Map) (*table.Table_Session, error) {
	t_id := table.Table_Session_Identifiers{Table_ID: form_values.Get("tableID"), Passcode: form_values.Get("passcode")}

	tsm := *table_sessions
	a, ok := tsm[t_id]
	if !ok {
		return nil, errors.New("table session does not exist")
	}

	a.Players = append(a.Players, *player.NewPlayer(form_values.Get("username")))

	return &a, nil
}
