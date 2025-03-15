package table

import (
	"dhdorr/story-point-poker/player"
)

type Table_Session_Identifiers struct {
	Table_ID string
	Passcode string
}

type Table_Map map[Table_Session_Identifiers]Table_Session

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
	new_table_settings := NewTableSettings(tsc.CL, tsc.NC, tsc.NR, tsc.TL, tsc.PM)
	new_player_arr := player.NewPlayerArr(tsc.PM)
	new_round_arr := NewPokerRoundArr(tsc.NR, tsc.PM)

	return &Table_Session{Table_ID: tsc.ID, Passcode: tsc.PC, Settings: *new_table_settings, Players: *new_player_arr, Rounds: *new_round_arr, Active_Round_ID: tsc.AR}
}

func (ts *Table_Session) AddPlayerToTableSession(un string) {
	ts.Players = append(ts.Players, *player.NewPlayer(un))
}

type Table_Session_Constructor struct {
	ID string
	PC string
	CL string
	UN string
	NC int
	NR int
	TL int
	PM int
	AR int
}
