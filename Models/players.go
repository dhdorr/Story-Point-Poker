package Models

type PLAYERS_DB struct {
	Players map[string]PLAYER
}

func (players_db PLAYERS_DB) RegisterPlayer(player PLAYER) {
	players_db.Players[player.Player_name] = player
}

type PLAYER struct {
	Player_name string
}

func CreatePlayer() PLAYER {
	player := new(PLAYER)
	return *player
}
