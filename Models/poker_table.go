package Models

type POKER_TABLE struct {
	Config          CONFIG_POKER_TABLE
	Players_DB      PLAYERS_DB
	Rounds_DB       ROUNDS_DB
	Admin           string
	CurrentRound    int
	PokerTableState int
}

func CreatePokerTable() POKER_TABLE {

	poker_table := POKER_TABLE{
		Config:          CONFIG_POKER_TABLE{},
		Players_DB:      PLAYERS_DB{Players: make(map[string]PLAYER)},
		Rounds_DB:       ROUNDS_DB{},
		Admin:           "",
		CurrentRound:    0,
		PokerTableState: 0,
	}

	return poker_table
}

func (poker_table POKER_TABLE) DestroyPokerTable() {}

func (poker_table POKER_TABLE) RegisterPlayer(player PLAYER) {
	poker_table.Players_DB.RegisterPlayer(player)
}

type CONFIG_POKER_TABLE struct {
	Table_name        string
	Table_passcode    string
	Table_max_players int
	Table_card_type   int
	Table_card_max    int
}

func CreateConfig() CONFIG_POKER_TABLE {
	config := new(CONFIG_POKER_TABLE)

	return *config
}

func (config CONFIG_POKER_TABLE) Generate_Key() string {
	return config.Table_name + ":" + config.Table_passcode
}
