package templates

import "dhdorr/story-point-poker/card"

type Waiting struct {
	MaxPlayers  int
	PlayerCount int
	Username    string
	TableID     string
	Passcode    string
}

type Game_Table struct {
	Cards    []card.Card
	Username string
	TableID  string
	Passcode string
}
