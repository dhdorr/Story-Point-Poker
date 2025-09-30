package Models

type CREATE_POKER_TABLE_REQUEST struct {
	Table_name                  string
	Table_passcode              string
	Player_name                 string
	Table_max_number_of_players int
	Table_max_number_of_rounds  int
	Table_time_to_live          int
}

func (req CREATE_POKER_TABLE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

type JOIN_POKER_TABLE_REQUEST struct {
	Table_name     string
	Table_passcode string
	Player_name    string
}

func (req JOIN_POKER_TABLE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
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

func (req CONFIGURE_ROUND_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

type TRANSITION_TABLE_STATE_REQUEST struct {
	Table_name     string
	Table_passcode string
}

func (req TRANSITION_TABLE_STATE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

type SUBMIT_VOTE_REQUEST struct {
	Value          int
	Voter_name     string
	Table_name     string
	Table_passcode string
}

func (req SUBMIT_VOTE_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}

type MOCK_REQUEST struct {
	Table_name     string
	Table_passcode string
}

func (req MOCK_REQUEST) GenerateKey() string {
	return req.Table_name + ":" + req.Table_passcode
}
