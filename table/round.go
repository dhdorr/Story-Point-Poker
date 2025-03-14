package table

import (
	"dhdorr/story-point-poker/player"
	"fmt"
	"time"
)

type RoundPhase int

// If the round is waiting for players or started, it is the active round
const (
	PhaseNotStarted RoundPhase = iota
	PhaseWaitingForPlayers
	PhaseStarted
	PhaseFinished
)

type Poker_Round struct {
	Start_Time time.Time
	Votes      map[*player.Player]int
	Phase      RoundPhase
}

func NewPokerRound(playerMax int) *Poker_Round {
	return &Poker_Round{Votes: make(map[*player.Player]int), Phase: PhaseNotStarted}
}

func NewPokerRoundArr(num_rounds, player_max int) *[]Poker_Round {
	arr := make([]Poker_Round, 0, num_rounds)
	for i := 0; i < num_rounds; i++ {
		npr := *NewPokerRound(player_max)
		fmt.Println(npr)
		arr = append(arr, npr)
	}
	return &arr
}

func (round *Poker_Round) SubmitVote(player *player.Player, vote int) {
	round.Votes[player] = vote
}
