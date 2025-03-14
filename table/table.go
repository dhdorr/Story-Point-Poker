package table

import (
	"dhdorr/story-point-poker/player"
	"net/url"
	"strconv"
)

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
	Table_ID        string
	Passcode        string
	Settings        Table_Settings
	Players         []player.Player
	Rounds          []Poker_Round
	Active_Round_ID int
}

func NewTableSessionConstructed(tsc Table_Session_Constructor) *Table_Session {
	new_table_settings := NewTableSettings(tsc.cl, tsc.nc, tsc.nr, tsc.tl, tsc.pm)
	new_player_arr := player.NewPlayerArr(tsc.pm, *player.NewPlayer(tsc.un))
	new_round_arr := NewPokerRoundArr(tsc.nr, tsc.pm)

	return &Table_Session{Table_ID: tsc.id, Passcode: tsc.pc, Settings: *new_table_settings, Players: *new_player_arr, Rounds: *new_round_arr, Active_Round_ID: tsc.ar}
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
	ar int
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

	tsc.ar = -1

	return NewTableSessionConstructed(tsc), nil
}
