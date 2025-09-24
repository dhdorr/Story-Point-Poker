package main

import (
	"fmt"

	"github.com/dhdorr/story-point-poker/Models"
)

type POKER_TABLE_DB struct {
	poker_tables map[string]Models.POKER_TABLE
}

func (poker_table_db POKER_TABLE_DB) CheckIfPokerTableExists(key string) bool {
	_, ok := poker_table_db.poker_tables[key]
	return ok
}

func (poker_table_db POKER_TABLE_DB) RegisterPokerTable(key string, poker_table Models.POKER_TABLE) {
	poker_table_db.poker_tables[key] = poker_table
}

func (poker_table_db POKER_TABLE_DB) ProcessCreateNewPokerTableRequest(req Models.CREATE_POKER_TABLE_REQUEST) {
	key := req.GenerateKey()
	if poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("A poker table already exists with the key: %s", key)
		return
	}

	poker_table := Models.InitializeNewPokerTableFromRequest(req)
	poker_table_db.RegisterPokerTable(key, poker_table)

	poker_table = poker_table_db.poker_tables[key]
	poker_table.PokerTableState = poker_table.TransitionPokerTableState()
	poker_table_db.poker_tables[key] = poker_table

	// Respond with round configuration page
}

func (poker_table_db POKER_TABLE_DB) ProcessConfigureRoundRequest(req Models.CONFIGURE_ROUND_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	poker_table := poker_table_db.poker_tables[key]
	round, err := poker_table.Rounds_DB.GetCurrentRound(poker_table.CurrentRound)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	round = round.ApplyConfigurationFromRequest(req)

	poker_table.Rounds_DB.Rounds[poker_table.CurrentRound] = round

}

func (poker_table_db POKER_TABLE_DB) ProcessJoinPokerTableRequest(req Models.JOIN_POKER_TABLE_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	poker_table := poker_table_db.poker_tables[key]

	player := Models.PLAYER{}
	poker_table.RegisterPlayer(player)

	poker_table_db.poker_tables[key] = poker_table
}

func (poker_table_db POKER_TABLE_DB) ProcessTransitionTableStateRequest(req Models.TRANSITION_TABLE_STATE_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	poker_table := poker_table_db.poker_tables[key]
	state := poker_table.TransitionPokerTableState()
	poker_table.PokerTableState = state
	poker_table_db.poker_tables[key] = poker_table
}

func main() {
	fmt.Println("Starting the server on :3000...")

	poker_table_db := new(POKER_TABLE_DB)
	poker_table_db.poker_tables = make(map[string]Models.POKER_TABLE)

	BeginTests(poker_table_db)
}

// TESTING AREA
type poker_table_test_interface interface {
	TestCreatePokerTable()
	TestConfigureFirstRound()
	TestConfigureSecondRound()
	TestTransitionPokerTableState()
	TestJoinPokerTable()
}

func BeginTests(pt poker_table_test_interface) {
	pt.TestCreatePokerTable()
	pt.TestConfigureFirstRound()
	// pt.TestTransitionPokerTableState()
	// pt.TestJoinPokerTable()
	// pt.TestTransitionPokerTableState()
	// pt.TestConfigureSecondRound()
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

func (poker_table_db POKER_TABLE_DB) TestConfigureFirstRound() {
	configure_round_request := Models.CONFIGURE_ROUND_REQUEST{
		Table_name:        "test",
		Table_passcode:    "test",
		Round_Title:       "title",
		Round_Description: "description",
		Card_Type:         0,
		Number_Of_Cards:   8,
		Time_Limit:        10,
	}

	poker_table_db.ProcessConfigureRoundRequest(configure_round_request)
	fmt.Printf("ROUNDS (%v): %v\n", len(poker_table_db.poker_tables["test:test"].Rounds_DB.Rounds), poker_table_db.poker_tables["test:test"].Rounds_DB)
}

func (poker_table_db POKER_TABLE_DB) TestConfigureSecondRound() {
	configure_round_request := Models.CONFIGURE_ROUND_REQUEST{
		Table_name:        "test",
		Table_passcode:    "test",
		Round_Title:       "title2",
		Round_Description: "description2",
		Card_Type:         0,
		Number_Of_Cards:   8,
		Time_Limit:        10,
	}

	poker_table_db.ProcessConfigureRoundRequest(configure_round_request)
	fmt.Printf("ROUNDS (%v): %v\n", len(poker_table_db.poker_tables["test:test"].Rounds_DB.Rounds), poker_table_db.poker_tables["test:test"].Rounds_DB)
}

func (poker_table_db POKER_TABLE_DB) TestJoinPokerTable() {
	join_poker_table_request := Models.JOIN_POKER_TABLE_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
		Player_name:    "apollo",
	}

	poker_table_db.ProcessJoinPokerTableRequest(join_poker_table_request)

	fmt.Printf("Poker Tables: %v\n", poker_table_db.poker_tables)
}

func (poker_table_db POKER_TABLE_DB) TestTransitionPokerTableState() {
	transition_poker_table_state_request := Models.TRANSITION_TABLE_STATE_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
	}

	poker_table_db.ProcessTransitionTableStateRequest(transition_poker_table_state_request)

	fmt.Printf("STATE: %v\n", poker_table_db.poker_tables["test:test"].PokerTableState)
}

// END TESTING
