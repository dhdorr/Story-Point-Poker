package player

type Player struct {
	Username string
	GUID     string
}

func GeneratePlayerGUID() string {
	return "temp-guid"
}

func NewPlayer(username string) *Player {
	return &Player{Username: username, GUID: GeneratePlayerGUID()}
}

func NewPlayerArr(player_max int, player Player) *[]Player {
	arr := make([]Player, 0, player_max)
	arr = append(arr, player)
	return &arr
}
