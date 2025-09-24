package Models

import (
	"fmt"
)

type ROUND_STATE int

const (
	StateRoundInitializing ROUND_STATE = iota
	StateRoundStarted
	StateRoundEnded
)

type ROUNDS_DB struct {
	Rounds []ROUND
}

type ROUND struct {
	Votes_DB    map[string]VOTE
	Config      ROUND_CONFIG
	Round_State ROUND_STATE
}

type ROUND_CONFIG struct {
	Round_Title       string
	Round_Description string
	Card_Type         int
	Number_Of_Cards   int
	Time_Limit        int
}

func (round ROUND) TransitionRoundState() ROUND_STATE {
	state := round.Round_State

	switch state {
	case StateRoundInitializing:
		return StateRoundStarted
	case StateRoundStarted:
		return StateRoundEnded
	default:
		panic(fmt.Errorf("unkown state: %v", state))
	}
}

func CreateRoundsDB(req CREATE_POKER_TABLE_REQUEST) ROUNDS_DB {
	rounds_db := ROUNDS_DB{Rounds: make([]ROUND, 0, 10)}
	return rounds_db
}

func CreateRound() ROUND {
	round := ROUND{}
	return round
}

func (rounds_db *ROUNDS_DB) RegisterRound(round ROUND) {
	rounds_db.Rounds = append(rounds_db.Rounds, round)
}

func (round ROUND) ApplyConfigurationFromRequest(req CONFIGURE_ROUND_REQUEST) ROUND {
	round.Config = ROUND_CONFIG{
		Round_Title:       req.Round_Title,
		Round_Description: req.Round_Description,
		Card_Type:         req.Card_Type,
		Number_Of_Cards:   req.Number_Of_Cards,
		Time_Limit:        req.Time_Limit,
	}

	return round
}

func (rounds_db ROUNDS_DB) GetCurrentRound(round_number int) (ROUND, error) {
	if len(rounds_db.Rounds) <= round_number || round_number < 0 {
		return ROUND{}, fmt.Errorf("index: %v is out of bounds", round_number)
	}
	round := rounds_db.Rounds[round_number]

	return round, nil
}

func (config ROUND_CONFIG) ConfigureRoundConfig(req CONFIGURE_ROUND_REQUEST) ROUND_CONFIG {
	config.Card_Type = req.Card_Type
	config.Number_Of_Cards = req.Number_Of_Cards
	config.Time_Limit = req.Time_Limit

	return config
}
