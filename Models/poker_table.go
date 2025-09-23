package Models

import "fmt"

type POKER_TABLE_STATE int

const (
	StateInitializing POKER_TABLE_STATE = iota
	StateWaitingForPlayers
	StateStarted
	StateEnded
)

type POKER_TABLE struct {
	Config          CONFIG_POKER_TABLE
	Players_DB      PLAYERS_DB
	Rounds_DB       ROUNDS_DB
	Admin           string
	CurrentRound    int
	PokerTableState POKER_TABLE_STATE
}

func (poker_table POKER_TABLE) TransitionPokerTableState(override_state int) POKER_TABLE_STATE {
	state := poker_table.PokerTableState

	if override_state != 0 {
		return POKER_TABLE_STATE(override_state)
	}

	switch state {
	case StateInitializing:
		return StateWaitingForPlayers
	case StateWaitingForPlayers:
		return StateStarted
	case StateStarted:
		return StateEnded
	default:
		panic(fmt.Errorf("unkown state: %v", state))
	}
}

func CreatePokerTable() POKER_TABLE {

	poker_table := POKER_TABLE{
		Config:          CreateConfig(),
		Players_DB:      PLAYERS_DB{Players: make(map[string]PLAYER)},
		Rounds_DB:       ROUNDS_DB{Rounds: make([]ROUND, 0, 10)},
		Admin:           "",
		CurrentRound:    0,
		PokerTableState: StateInitializing,
	}

	return poker_table
}

func (poker_table POKER_TABLE) DestroyPokerTable() {}

func (poker_table POKER_TABLE) RegisterPlayer(player PLAYER) {
	poker_table.Players_DB.RegisterPlayer(player)
}

type CONFIG_POKER_TABLE struct {
	Table_name                  string
	Table_passcode              string
	Table_max_number_of_players int
	Table_max_number_of_rounds  int
	Table_time_to_live          int
}

func CreateConfig() CONFIG_POKER_TABLE {
	config := new(CONFIG_POKER_TABLE)

	return *config
}

func (config CONFIG_POKER_TABLE) Generate_Key() string {
	return config.Table_name + ":" + config.Table_passcode
}
