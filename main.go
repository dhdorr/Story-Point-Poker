package main

import (
	"dhdorr/story-point-poker/managers"
	"dhdorr/story-point-poker/table"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to Story Point Poker 2")

	// Global table manager that keeps every active table session in memory
	tm := managers.Table_Manager{Table_Sessions_M: make(map[table.Table_Session_Identifiers]table.Table_Session)}

	// Test
	tm.HandleTest()
	tm.PrintTables()

	http.Handle("/", http.FileServer(http.Dir(".")))

	// Serve static pages
	fs_pages := http.FileServer(http.Dir("pages"))
	http.HandleFunc("GET /pages/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getting some page...")
		http.StripPrefix("/pages/", fs_pages).ServeHTTP(w, r)
	})

	// Custom API Requests
	http.HandleFunc("POST /joinTable", tm.HandleJoin)

	http.HandleFunc("POST /createTable", tm.HandleCreate)
	// *******************

	// serve css and js, from html pages
	fs_static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
