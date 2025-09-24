package Models

import "fmt"

type POKER_TABLE_STATE int

const (
	StateInitializing POKER_TABLE_STATE = iota
	StateStarted
	StateEnded
)

type POKER_TABLE struct {
	Config          CONFIG_POKER_TABLE
	Players_DB      PLAYERS_DB
	Rounds_DB       ROUNDS_DB
	Admin           PLAYER
	CurrentRound    int
	PokerTableState POKER_TABLE_STATE
}

type CONFIG_POKER_TABLE struct {
	Table_name                  string
	Table_passcode              string
	Table_max_number_of_players int
	Table_max_number_of_rounds  int
	Table_time_to_live          int
}

func (poker_table *POKER_TABLE) TransitionPokerTableState() POKER_TABLE_STATE {
	state := poker_table.PokerTableState

	switch state {
	case StateInitializing:
		round := CreateRound()
		poker_table.Rounds_DB.RegisterRound(round)
		return StateStarted
	case StateStarted:
		return StateEnded
	default:
		panic(fmt.Errorf("unkown state: %v", state))
	}
}

func InitializeNewPokerTableFromRequest(req CREATE_POKER_TABLE_REQUEST) POKER_TABLE {

	poker_table := POKER_TABLE{
		Config:          createConfig(req),
		Players_DB:      CreatePlayersDB(req),
		Rounds_DB:       CreateRoundsDB(req),
		Admin:           PLAYER{},
		CurrentRound:    0,
		PokerTableState: StateInitializing,
	}

	player := CreatePlayer(req.Player_name)
	poker_table.Players_DB.RegisterPlayer(player)
	poker_table.Admin = player

	return poker_table
}

func createConfig(req CREATE_POKER_TABLE_REQUEST) CONFIG_POKER_TABLE {
	config := CONFIG_POKER_TABLE{
		Table_name:                  req.Table_name,
		Table_passcode:              req.Table_passcode,
		Table_max_number_of_players: req.Table_max_number_of_players,
		Table_max_number_of_rounds:  req.Table_max_number_of_rounds,
		Table_time_to_live:          req.Table_time_to_live,
	}
	return config
}

func (poker_table POKER_TABLE) DestroyPokerTable() {}

func (poker_table POKER_TABLE) RegisterPlayer(player PLAYER) {
	poker_table.Players_DB.RegisterPlayer(player)
}

func (config CONFIG_POKER_TABLE) GenerateKey() string {
	return config.Table_name + ":" + config.Table_passcode
}
