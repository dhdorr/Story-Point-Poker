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
	IsDone   bool
}

type Game_Update struct {
	IsDone bool
}

type Results_Card struct {
	Value int
	Votes int
}

type Round_Results struct {
	Cards    []Results_Card
	Style    template.CSS
	Username string
	TableID  string
	Passcode string
	IsValid  bool
}
