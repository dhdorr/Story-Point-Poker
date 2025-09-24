package Models

type PLAYERS_DB struct {
	Players []PLAYER
}

func CreatePlayersDB(req CREATE_POKER_TABLE_REQUEST) PLAYERS_DB {
	players_db := PLAYERS_DB{Players: make([]PLAYER, 0, 10)}
	return players_db
}

func (players_db *PLAYERS_DB) RegisterPlayer(player PLAYER) {
	players_db.Players = append(players_db.Players, player)
}

type PLAYER struct {
	Player_name string
}

func CreatePlayer(name string) PLAYER {
	return PLAYER{Player_name: name}
}

func (players_db PLAYERS_DB) CreateAndRegisterPlayer(name string) PLAYER {
	player := CreatePlayer(name)
	players_db.RegisterPlayer(player)

	return player
}
