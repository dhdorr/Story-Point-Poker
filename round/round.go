package round

import (
	"dhdorr/story-point-poker/card"
	"time"
)

type RoundPhase int

const (
	PhaseWaiting RoundPhase = iota
	PhaseStarted
	PhaseResults
	PhaseFinished
)

type Round struct {
	StartTime time.Time
	EndTime   time.Time
	Phase     RoundPhase
	Topic     string
	Votes     []card.Vote
}
