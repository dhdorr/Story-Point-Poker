package main

import (
	"dhdorr/story-point-poker/handlers"
	"dhdorr/story-point-poker/table_manager"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to Story Point Poker 3")

	tm := table_manager.Table_Manager{TableMap: make(table_manager.Table_Map)}

	// Serve home page
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Custom API Requests
	http.HandleFunc("GET /home", handlers.HandleHome)
	http.HandleFunc("GET /join-table-menu", handlers.HandleJoinTableMenu)
	http.HandleFunc("GET /create-table-menu", handlers.HandleCreateTableMenu)
	http.HandleFunc("POST /join-table", tm.HandleJoinTable)
	http.HandleFunc("POST /create-table", tm.HandleCreateTable)

	// *******************

	// serve css and js, from html pages
	fs_static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
