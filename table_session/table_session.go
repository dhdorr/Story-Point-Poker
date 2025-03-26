package table_session

import (
	"dhdorr/story-point-poker/card"
	"dhdorr/story-point-poker/player"
	"dhdorr/story-point-poker/round"
	"time"
)

type TableState int // accept players, closed to new players, game is done

const (
	StateOpen TableState = iota
	StateClosed
	StateDone
)

type Settings struct {
	CardLayout     string
	NumCards       int
	NumRounds      int
	MaxPlayers     int
	RoundTimeLimit int
	TableTimeLimit int
}

type Table struct {
	TableID       string
	Passcode      string
	ActiveRoundID int
	State         TableState
	Players       []player.Player
	Rounds        []round.Round
	Cards         []card.Card
	Admin         player.Player
	StartTime     time.Time
	EndTime       time.Time
	Settings      Settings
}
