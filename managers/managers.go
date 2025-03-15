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
	ts, err := handlers.HandleJoin(w, r, &tm.Table_Sessions_M)
	if err != nil {
		fmt.Fprintf(w, "Failed to join table session: %v", err)
		return
	}

	filename := "T-poker-table.html"
	handlers.RenderTemplate(w, filename, *ts)
}

func (tm *Table_Manager) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ts, err := handlers.HandleCreate(w, r)
	if err != nil {
		fmt.Fprintf(w, "Failed to create a new table session: %v", err)
		return
	}

	tm.Table_Sessions_M[table.Table_Session_Identifiers{Table_ID: ts.Table_ID, Passcode: ts.Passcode}] = *ts

	filename := "T-poker-table.html"
	handlers.RenderTemplate(w, filename, *ts)
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
}

func HandleTestCreate(form_values url.Values) (*table.Table_Session, error) {
	ts, err := handlers.GenerateTableSession(form_values)
	if err != nil {
		return nil, err
	}

	fmt.Printf("session made: %v \n", ts)
	return ts, nil
}
