package templates

import "dhdorr/story-point-poker/card"

type Waiting struct {
	MaxPlayers  int
	PlayerCount int
}

type Game_Table struct {
	Cards []card.Card
}
