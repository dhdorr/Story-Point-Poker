package player

type Player struct {
	Username string
	GUID     string
	IsAdmin  bool
}

func GeneratePlayerGUID() string {
	return "temp-guid"
}

func NewPlayer(username string, isAdmin bool) *Player {
	return &Player{Username: username, GUID: GeneratePlayerGUID(), IsAdmin: isAdmin}
}

func NewPlayerArr(player_max int) *[]Player {
	arr := make([]Player, 0, player_max)
	return &arr
}
