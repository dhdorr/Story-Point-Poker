package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Player struct {
	Username string
	GUID     string
}

func GeneratePlayerGUID() string {
	return "temp-guid"
}

func NewPlayer(username string) *Player {
	return &Player{Username: username, GUID: GeneratePlayerGUID()}
}

func NewPlayerArr(player_max int, player Player) *[]Player {
	arr := make([]Player, 0, player_max)
	arr = append(arr, player)
	return &arr
}

type Table_Settings struct {
	Card_Layout      string
	Number_Of_Cards  int
	Number_Of_Rounds int
	Round_Time_Limit int
	Player_Max       int
}

func NewTableSettings(card_layout string, num_cards, num_rounds, round_time_limit, player_max int) *Table_Settings {
	return &Table_Settings{Card_Layout: card_layout, Number_Of_Cards: num_cards, Number_Of_Rounds: num_rounds, Round_Time_Limit: round_time_limit, Player_Max: player_max}
}

type Table_Session struct {
	Table_ID string
	Passcode string
	Settings Table_Settings
	Players  []Player
}

func NewTableSession(table_id, passcode, card_layout, username string, num_cards, num_rounds, round_time_limit, player_max int) *Table_Session {
	return &Table_Session{Table_ID: table_id, Passcode: passcode, Settings: *NewTableSettings(card_layout, num_cards, num_rounds, round_time_limit, player_max), Players: *NewPlayerArr(player_max, *NewPlayer(username))}
}

func NewTableSessionConstructed(tsc Table_Session_Constructor) *Table_Session {
	return &Table_Session{Table_ID: tsc.id, Passcode: tsc.pc, Settings: *NewTableSettings(tsc.cl, tsc.nc, tsc.nr, tsc.tl, tsc.pm), Players: *NewPlayerArr(tsc.pm, *NewPlayer(tsc.un))}
}

type Table_Session_Constructor struct {
	id string
	pc string
	cl string
	un string
	nc int
	nr int
	tl int
	pm int
}

func GenerateTableSession(form_values url.Values) (*Table_Session, error) {
	tsc := Table_Session_Constructor{}
	tsc.id = form_values.Get("tableID")
	tsc.pc = form_values.Get("passcode")
	tsc.cl = form_values.Get("cardLayout")
	tsc.un = form_values.Get("username")

	nc, err := strconv.Atoi(form_values.Get("numCards"))
	if err != nil {
		return nil, err
	}
	tsc.nc = nc

	nr, err := strconv.Atoi(form_values.Get("numRounds"))
	if err != nil {
		return nil, err
	}
	tsc.nr = nr

	tl, err := strconv.Atoi(form_values.Get("roundTimeLimit"))
	if err != nil {
		return nil, err
	}
	tsc.tl = tl

	pm, err := strconv.Atoi(form_values.Get("playerMax"))
	if err != nil {
		return nil, err
	}
	tsc.pm = pm

	return NewTableSessionConstructed(tsc), nil
}

func test_GTS() {
	test, err := GenerateTableSession(url.Values{
		"tableID":        []string{"test"},
		"passcode":       []string{"test"},
		"cardLayout":     []string{"test"},
		"username":       []string{"test"},
		"numCards":       []string{"1"},
		"numRounds":      []string{"2"},
		"roundTimeLimit": []string{"3"},
		"playerMax":      []string{"4"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("TEST: %v \n", test)
}

func main() {
	fmt.Println("Welcome to Story Point Poker 2")

	// TESTING
	test_GTS()
	// END TESTING

	http.Handle("/", http.FileServer(http.Dir(".")))

	fs_pages := http.FileServer(http.Dir("pages"))
	http.HandleFunc("GET /pages/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/pages/", fs_pages).ServeHTTP(w, r)
	})

	// // serve dynamic html templates
	// fs_templates := http.FileServer(http.Dir("templates"))
	// http.HandleFunc("GET /templates/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.RequestURI)
	// 	http.StripPrefix("/templates/", fs_templates).ServeHTTP(w, r)
	// })

	http.HandleFunc("POST /joinTable", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Printf("joining table: %v \n", r.Form)
		tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("POST /createTable", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Printf("creating table: %v \n", r.Form)
		tml, err := GenerateTableSession(r.Form)
		if err != nil {
			fmt.Fprintf(w, "%s", err.Error())
			return
		}

		fmt.Printf("session made: %v \n", tml)

		tmpl := template.Must(template.ParseFiles("templates/T-poker-table.html"))
		tmpl.Execute(w, nil)
	})

	// serve css and js, from html pages
	fs_static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
