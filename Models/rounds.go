package Models

type ROUND_STATE int

const (
	StateInitializingRound ROUND_STATE = iota
	StateStartedRound
	StateEndedRound
)

type ROUNDS_DB struct {
	Rounds []ROUND
}

type ROUND struct {
	Votes_DB          map[string]VOTE
	Round_Title       string
	Round_Description string
	Config            ROUND_CONFIG
	Round_State       ROUND_STATE
}

type ROUND_CONFIG struct {
	Card_Type       int
	Number_Of_Cards int
	Time_Limit      int
}

func CreateRound() ROUND {
	round := new(ROUND)
	return *round
}

func (round ROUND) ConfigureRound(req CONFIGURE_ROUND_REQUEST) ROUND {
	round.Config = round.Config.ConfigureRoundConfig(req)
	round.Round_Description = req.Round_Description
	round.Round_Title = req.Round_Title

	return round
}

func (config ROUND_CONFIG) ConfigureRoundConfig(req CONFIGURE_ROUND_REQUEST) ROUND_CONFIG {
	config.Card_Type = req.Card_Type
	config.Number_Of_Cards = req.Number_Of_Cards
	config.Time_Limit = req.Time_Limit

	return config
}
