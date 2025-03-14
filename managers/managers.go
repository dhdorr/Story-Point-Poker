package managers

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/table"
	"fmt"
	"net/http"
)

type Table_Manager struct {
	Table_Sessions []table.Table_Session
}

func (tm *Table_Manager) HandleJoin(w http.ResponseWriter, r *http.Request) {
	handlers.HandleJoin(w, r)
}

func (tm *Table_Manager) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ts, err := handlers.HandleCreate(w, r)
	if err != nil {
		fmt.Fprintf(w, "Failed to create a new table session: %s", err.Error())
		return
	}
	tm.Table_Sessions = append(tm.Table_Sessions, *ts)

	handlers.RenderTemplate(w, *ts)
}
