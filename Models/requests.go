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

func (c_r CREATE_POKER_TABLE_REQUEST) GenerateKey() string {
	return c_r.Table_name + ":" + c_r.Table_passcode
}

func (c_r CREATE_POKER_TABLE_REQUEST) GenerateConfig() CONFIG_POKER_TABLE {
	config := CreateConfig()
	config.Table_name = c_r.Table_name
	config.Table_passcode = c_r.Table_passcode

	return config
}

func (c_r CREATE_POKER_TABLE_REQUEST) GeneratePlayer() PLAYER {
	player := CreatePlayer()

	player.Player_name = c_r.Player_name
	return player
}

func (j_r JOIN_POKER_TABLE_REQUEST) GenerateKey() string {
	return j_r.Table_name + ":" + j_r.Table_passcode
}

func (j_r JOIN_POKER_TABLE_REQUEST) GenerateConfig() CONFIG_POKER_TABLE {
	config := CreateConfig()
	config.Table_name = j_r.Table_name
	config.Table_passcode = j_r.Table_passcode

	return config
}

func (j_r JOIN_POKER_TABLE_REQUEST) GeneratePlayer() PLAYER {
	player := CreatePlayer()

	player.Player_name = j_r.Player_name
	return player
}
