package managers

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/table"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Table_Manager struct {
	Table_Sessions_M table.Table_Map
}

func (tm *Table_Manager) PrintTables() {
	fmt.Printf("tm.Table_Sessions_M: %v\n", tm.Table_Sessions_M)
}

func (tm *Table_Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	t_id := table.Table_Session_Identifiers{Table_ID: r.FormValue("tableID"), Passcode: r.FormValue("passcode")}
	un := r.FormValue("username")

	ts, ok := tm.Table_Sessions_M[t_id]
	if !ok {
		fmt.Fprintf(w, "Unable to find table: %v", t_id.Table_ID)
		return
	}

	filename := "T-poker-table.html"
	// TODO:
	// Check if session is waiting for players or active or closed!
	// if ts.IsOpen ...
	if ts.Rounds[ts.Active_Round_ID].Phase == table.PhaseWaitingForPlayers {
		filename = "T-waiting.html"
	}
	// ^^^

	ts.AddPlayerToTableSession(un)

	tm.Table_Sessions_M[t_id] = ts

	fmt.Printf("Session Joined: %v\n", ts.Players)

	tm.PrintTables()

	handlers.RenderTemplate(w, filename, ts)
}

func (tm *Table_Manager) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ts, err := handlers.HandleCreate(w, r)
	if err != nil {
		fmt.Fprintf(w, "Failed to create a new table session: %v", err)
		return
	}

	tm.AddNewTableSession(table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}, ts)
	tm.PrintTables()

	filename := "T-waiting.html"
	// w.Header().Add("tableID", ts.Table_ID)
	handlers.RenderTemplate(w, filename, *ts)
}

func (tm *Table_Manager) HandleCheckForNewPlayers(w http.ResponseWriter, r *http.Request) {
	t_id := r.URL.Query().Get("tableID")
	pc := r.URL.Query().Get("passcode")
	// un := r.URL.Query().Get("username")

	ts := tm.Table_Sessions_M[*table.NewTableSessionIdentifier(t_id, pc)]
	ts.PrintTableSessionPlayers()

	filename := "T-players.html"
	handlers.RenderTemplate(w, filename, ts)
}

func (tm *Table_Manager) HandleStart(w http.ResponseWriter, r *http.Request) {
	keys := make([]table.Table_Session_Identifiers, 0, len(tm.Table_Sessions_M))
	for t := range tm.Table_Sessions_M {
		keys = append(keys, t)
	}

	// fmt.Println(r.URL.Query())
	t := tm.Table_Sessions_M[keys[0]]
	t.Active_Round_ID, _ = strconv.Atoi(r.URL.Query().Get("activeRound"))
	tm.Table_Sessions_M[keys[0]] = t

	filename := "T-poker-table.html"
	handlers.RenderTemplate(w, filename, tm.Table_Sessions_M[keys[0]])
}

func (tm *Table_Manager) AddNewTableSession(t_id table.Table_Session_Identifiers, ts *table.Table_Session) {
	tm.Table_Sessions_M[t_id] = *ts
}

// TESTING

func (tm *Table_Manager) HandleTest() {
	test := url.Values{
		"tableID":        []string{"test_id"},
		"passcode":       []string{"test_pc"},
		"cardLayout":     []string{"seq"},
		"username":       []string{"test_un"},
		"numCards":       []string{"6"},
		"numRounds":      []string{"1"},
		"roundTimeLimit": []string{"30"},
		"playerMax":      []string{"10"},
	}
	ts, err := HandleTestCreate(test)
	if err != nil {
		fmt.Println("Failed to create a test table session")
		return
	}

	tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}] = *ts

	ts_j := tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}]

	tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}] = *HandleTestJoin(test, &ts_j)
}

func HandleTestCreate(form_values url.Values) (*table.Table_Session, error) {
	ts, err := handlers.GenerateTableSession(form_values)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func HandleTestJoin(form_values url.Values, tmp *table.Table_Session) *table.Table_Session {
	un := form_values.Get("username")
	tmp.AddPlayerToTableSession(un)
	return tmp
}
