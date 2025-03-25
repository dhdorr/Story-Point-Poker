package templates

import (
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/table"
)

type Gen_Test_A struct {
	Value int
}

type Gen_Test_B struct {
	Data string
}

type Gen_Test_Interface interface {
	Gen_Test_A | Gen_Test_B
}

type Waiting_For_Players struct {
	Table_ID    string
	Passcode    string
	PlayerCount int
	MaxPlayers  int
	Username    string
	GUID        string
	IsAdmin     bool
	IsReady     int
	TimeLimit   int
}

type Player_Count struct {
	PlayerCount int
	IsReady     int
}

type Game_Table struct {
	Cards    []table.Card
	Players  []player.Player
	Table_ID string
	Passcode string
	IsAdmin  bool
	Username string
}

type Results struct {
	Cards       []table.Results_Card
	ActiveRound int
	NextRound   int
	IsAdmin     bool
	Table_ID    string
	Passcode    string
	Username    string
}

type Should_Change_Round struct {
	ChangeRound bool
}

type Gen_Table_Session_Interface interface {
	table.Table_Session | table.Table_Session_Constructor | Gen_Test_B | Gen_Test_A | table.Poker_Round | Player_Count | Waiting_For_Players | Results | Game_Table | Should_Change_Round
}
