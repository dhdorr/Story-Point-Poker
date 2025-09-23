package main

import (
	"fmt"

	"github.com/dhdorr/story-point-poker/Models"
)

type POKER_TABLE_DB struct {
	poker_tables map[string]Models.POKER_TABLE
}

func (poker_table_db POKER_TABLE_DB) Generate_Key(table_name, table_passcode string) string {
	return table_name + ":" + table_passcode
}

func (poker_table_db POKER_TABLE_DB) CheckIfPokerTableExists(key string) bool {
	_, ok := poker_table_db.poker_tables[key]
	return ok
}

func CreateNewPokerTable() Models.POKER_TABLE {
	poker_table := Models.CreatePokerTable()

	return poker_table
}

func (poker_table_db POKER_TABLE_DB) RegisterPokerTable(key string, poker_table Models.POKER_TABLE) {
	poker_table_db.poker_tables[key] = poker_table
}

func CreateNewPlayer(req Models.PLAYER_REQUEST_INTERFACE) Models.PLAYER {
	return req.GeneratePlayer()
}

func (poker_table_db POKER_TABLE_DB) ProcessCreateNewPokerTableRequest(req Models.PLAYER_REQUEST_INTERFACE) {
	key := req.GenerateKey()
	if poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("A poker table already exists with the key: %s", key)
		return
	}

	poker_table := CreateNewPokerTable()
	config := req.GenerateConfig()
	poker_table.Config = config

	poker_table_db.RegisterPokerTable(key, poker_table)

	player := CreateNewPlayer(req)
	poker_table.RegisterPlayer(player)

	poker_table_db.poker_tables[key] = poker_table
}

func (poker_table_db POKER_TABLE_DB) ProcessJoinPokerTableRequest(req Models.CREATE_POKER_TABLE_REQUEST) {
	key := poker_table_db.Generate_Key(req.Table_name, req.Table_passcode)
	ok := poker_table_db.CheckIfPokerTableExists(key)
	if !ok {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	player := CreateNewPlayer(req)
	poker_table := poker_table_db.poker_tables[req.GenerateKey()]
	poker_table.RegisterPlayer(player)

	poker_table_db.poker_tables[key] = poker_table
}

type poker_table_interface interface {
	TestCreatePokerTable()
	TestJoinPokerTable()
}

func main() {
	fmt.Println("Starting the server on :3000...")

	poker_table_db := new(POKER_TABLE_DB)
	poker_table_db.poker_tables = make(map[string]Models.POKER_TABLE)

	BeginTests(poker_table_db)
}

// TESTING AREA
func BeginTests(pt poker_table_interface) {
	pt.TestCreatePokerTable()
	pt.TestJoinPokerTable()
}

func (poker_table_db POKER_TABLE_DB) TestCreatePokerTable() {
	create_poker_table_request := Models.CREATE_POKER_TABLE_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
		Player_name:    "artemis",
	}

	poker_table_db.ProcessCreateNewPokerTableRequest(create_poker_table_request)

	fmt.Printf("Poker Tables: %v\n", poker_table_db.poker_tables)
}

func (poker_table_db POKER_TABLE_DB) TestJoinPokerTable() {
	join_poker_table_request := Models.CREATE_POKER_TABLE_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
		Player_name:    "apollo",
	}

	poker_table_db.ProcessJoinPokerTableRequest(join_poker_table_request)

	fmt.Printf("Poker Tables: %v\n", poker_table_db.poker_tables)
}

// END TESTING
