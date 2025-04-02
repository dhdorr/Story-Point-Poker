package templates

import (
	"dhdorr/story-point-poker/card"
	"html/template"
)

type Waiting struct {
	MaxPlayers  int
	PlayerCount int
	Username    string
	TableID     string
	Passcode    string
	Ready       bool
	Style       template.CSS
}

type Game_Table struct {
	Cards    []card.Card
	Username string
	TableID  string
	Passcode string
}
