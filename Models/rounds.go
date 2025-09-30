package Models

import (
	"fmt"
)

type ROUND_STATE int

const (
	StateRoundInitializing ROUND_STATE = iota
	StateRoundStarted
	StateRoundResults
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

type ROUND_RESULTS struct {
	Round_Title string
	Votes       []VOTE
	Most_Voted  []VOTE
}

func (round ROUND) TransitionRoundState() ROUND_STATE {
	state := round.Round_State

	switch state {
	case StateRoundInitializing:
		return StateRoundStarted
	case StateRoundStarted:
		return StateRoundResults
	case StateRoundResults:
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

func (round *ROUND) ApplyConfigurationFromRequest(req CONFIGURE_ROUND_REQUEST) {
	round.Config = ROUND_CONFIG{
		Round_Title:       req.Round_Title,
		Round_Description: req.Round_Description,
		Card_Type:         req.Card_Type,
		Number_Of_Cards:   req.Number_Of_Cards,
		Time_Limit:        req.Time_Limit,
	}

	votes_db := make(map[string]VOTE)
	round.Votes_DB = votes_db
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

func (round *ROUND) CastVote(req SUBMIT_VOTE_REQUEST) {
	votes_db := round.Votes_DB

	vote := VOTE{}
	vote.Value = req.Value
	votes_db[req.Voter_name] = vote

	round.Votes_DB = votes_db
}

func (round ROUND) GenerateRoundResults() ROUND_RESULTS {
	votes := round.Votes_DB

	result := ROUND_RESULTS{}
	most_voted_map := make(map[VOTE]int)
	most_voted := VOTE{}
	vote_count := 0
	votes_list := make([]VOTE, 0, 10)

	for _, v := range votes {
		_, ok := most_voted_map[v]
		if ok {
			most_voted_map[v] += 1
		} else {
			most_voted_map[v] = 1
			votes_list = append(votes_list, v)
		}
	}

	for mv, count := range most_voted_map {
		if count >= vote_count {
			most_voted = mv
			vote_count = count
		}
	}

	result.Most_Voted = append(result.Most_Voted, most_voted)
	result.Round_Title = round.Config.Round_Title
	result.Votes = votes_list

	return result
}
