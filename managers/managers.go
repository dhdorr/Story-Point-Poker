package managers

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/table"
	"fmt"
	"net/http"
	"net/url"
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
	ts.AddPlayerToTableSession(un)

	fmt.Printf("Session Joined: %v", ts.Players)

	filename := "T-poker-table.html"
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

	filename := "T-poker-table.html"
	handlers.RenderTemplate(w, filename, *ts)
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

	fmt.Printf("session finalized: %v \n", tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}])
}

func HandleTestCreate(form_values url.Values) (*table.Table_Session, error) {
	fmt.Printf("creating table: %v \n", form_values)
	ts, err := handlers.GenerateTableSession(form_values)
	if err != nil {
		return nil, err
	}
	// ts.AddPlayerToTableSession(form_values.Get("username"))

	fmt.Printf("session made: %v \n", ts)
	return ts, nil
}

func HandleTestJoin(form_values url.Values, tmp *table.Table_Session) *table.Table_Session {
	un := form_values.Get("username")
	tmp.AddPlayerToTableSession(un)
	return tmp
}
