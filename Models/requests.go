package Models

type CREATE_POKER_TABLE_REQUEST struct {
	Table_name     string
	Table_passcode string
	Player_name    string
}

type JOIN_POKER_TABLE_REQUEST struct {
	Table_name     string
	Table_passcode string
	Player_name    string
}

type PLAYER_REQUEST_INTERFACE interface {
	GenerateKey() string
	GenerateConfig() CONFIG_POKER_TABLE
	GeneratePlayer() PLAYER
}

func (req CREATE_POKER_TABLE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

func (req CREATE_POKER_TABLE_REQUEST) GenerateConfig() CONFIG_POKER_TABLE {
	config := CreateConfig()
	config.Table_name = req.Table_name
	config.Table_passcode = req.Table_passcode

	return config
}

func (req CREATE_POKER_TABLE_REQUEST) GeneratePlayer() PLAYER {
	player := CreatePlayer()

	player.Player_name = req.Player_name
	return player
}

func (req JOIN_POKER_TABLE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

func (req JOIN_POKER_TABLE_REQUEST) GenerateConfig() CONFIG_POKER_TABLE {
	config := CreateConfig()
	config.Table_name = req.Table_name
	config.Table_passcode = req.Table_passcode

	return config
}

func (req JOIN_POKER_TABLE_REQUEST) GeneratePlayer() PLAYER {
	player := CreatePlayer()

	player.Player_name = req.Player_name
	return player
}
