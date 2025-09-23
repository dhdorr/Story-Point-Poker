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

type CONFIGURE_ROUND_REQUEST struct {
	Table_name        string
	Table_passcode    string
	Round_Title       string
	Round_Description string
	Card_Type         int
	Number_Of_Cards   int
	Time_Limit        int
}

type TRANSITION_TABLE_STATE_REQUEST struct {
	Table_name     string
	Table_passcode string
	Override_State int
}

type TRANSITION_TABLE_STATE_INTERFACE interface {
	GenerateKey() string
	GetData() TRANSITION_TABLE_STATE_REQUEST
}

func (req TRANSITION_TABLE_STATE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

func (req TRANSITION_TABLE_STATE_REQUEST) GetData() TRANSITION_TABLE_STATE_REQUEST {
	return req
}

type PLAYER_REQUEST_INTERFACE interface {
	GenerateKey() string
}

type CONFIGURE_ROUND_INTERFACE interface {
	GenerateKey() string
	GetData() CONFIGURE_ROUND_REQUEST
}

type CREATE_PLAYER_REQUEST_INTERFACE interface {
	GenerateKey() string
	GeneratePlayer() PLAYER
}

type CREATE_POKER_TABLE_REQUEST_INTERFACE interface {
	GenerateKey() string
	GeneratePlayer() PLAYER
	GenerateConfig() CONFIG_POKER_TABLE
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

func (req CONFIGURE_ROUND_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

func (req CONFIGURE_ROUND_REQUEST) GetData() CONFIGURE_ROUND_REQUEST {
	return req
}

type CREATE_OR_JOIN_REQUEST_INTERFACE interface {
	CREATE_POKER_TABLE_REQUEST | JOIN_POKER_TABLE_REQUEST
}
