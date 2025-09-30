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

	round.ApplyConfigurationFromRequest(req)
	round.Round_State = round.TransitionRoundState()

	poker_table.Rounds_DB.Rounds[poker_table.CurrentRound] = round

	// Respond with game table page
}

func (poker_table_db *POKER_TABLE_DB) ProcessSubmitVoteRequest(req Models.SUBMIT_VOTE_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	// TODO: do
	poker_table := poker_table_db.poker_tables[key]
	round, err := poker_table.Rounds_DB.GetCurrentRound(poker_table.CurrentRound)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	round.CastVote(req)
	fmt.Printf("Votes: %v\n", round.Votes_DB)
	poker_table.Rounds_DB.Rounds[0] = round
	poker_table_db.poker_tables[key] = poker_table
}

func (poker_table_db POKER_TABLE_DB) ProcessShowRoundResultsRequest(req Models.MOCK_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	// TODO: do
	poker_table := poker_table_db.poker_tables[key]
	round, err := poker_table.Rounds_DB.GetCurrentRound(poker_table.CurrentRound)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	results := round.GenerateRoundResults()
	fmt.Printf("RESULTS: %v", results)
}

func (poker_table_db POKER_TABLE_DB) ProcessProceedToNextRoundRequest(req Models.MOCK_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	// TODO: do
	poker_table := poker_table_db.poker_tables[key]
	round := Models.CreateRound()
	poker_table.Rounds_DB.RegisterRound(round)
	poker_table.CurrentRound += 1

	poker_table_db.poker_tables[key] = poker_table
	fmt.Printf("Rounds Count: %v \n", len(poker_table_db.poker_tables[key].Rounds_DB.Rounds))
}

func (poker_table_db POKER_TABLE_DB) ProcessJoinPokerTableRequest(req Models.JOIN_POKER_TABLE_REQUEST) {
	key := req.GenerateKey()
	if !poker_table_db.CheckIfPokerTableExists(key) {
		fmt.Printf("No poker table with key: %s exists!", key)
		return
	}

	poker_table := poker_table_db.poker_tables[key]
	poker_table.InitializeNewPlayerFromRequest(req)

	poker_table_db.poker_tables[key] = poker_table

	// Respond with game table page
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
	TestJoinPokerTable()
	TestSubmitVote()
	TestGetRoundResults()
	TestProceedNextRound()
}

func BeginTests(pt poker_table_test_interface) {
	pt.TestCreatePokerTable()
	pt.TestConfigureFirstRound()
	pt.TestJoinPokerTable()
	pt.TestSubmitVote()
	pt.TestGetRoundResults()
	pt.TestProceedNextRound()
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

func (poker_table_db POKER_TABLE_DB) TestJoinPokerTable() {
	join_poker_table_request := Models.JOIN_POKER_TABLE_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
		Player_name:    "apollo",
	}

	poker_table_db.ProcessJoinPokerTableRequest(join_poker_table_request)

	fmt.Printf("Players: %v\n", poker_table_db.poker_tables["test:test"].Players_DB.Players)
}

func (poker_table_db POKER_TABLE_DB) TestSubmitVote() {
	submit_vote_request := Models.SUBMIT_VOTE_REQUEST{
		Value:          3,
		Voter_name:     "artemis",
		Table_name:     "test",
		Table_passcode: "test",
	}

	poker_table_db.ProcessSubmitVoteRequest(submit_vote_request)

	submit_vote_request2 := Models.SUBMIT_VOTE_REQUEST{
		Value:          5,
		Voter_name:     "apollo",
		Table_name:     "test",
		Table_passcode: "test",
	}

	poker_table_db.ProcessSubmitVoteRequest(submit_vote_request2)

	submit_vote_request3 := Models.SUBMIT_VOTE_REQUEST{
		Value:          3,
		Voter_name:     "derek",
		Table_name:     "test",
		Table_passcode: "test",
	}

	poker_table_db.ProcessSubmitVoteRequest(submit_vote_request3)
	fmt.Printf("VOTES: %v \n", poker_table_db.poker_tables["test:test"].Rounds_DB.Rounds[0].Votes_DB)
}

func (poker_table_db POKER_TABLE_DB) TestGetRoundResults() {
	mock := Models.MOCK_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
	}
	poker_table_db.ProcessShowRoundResultsRequest(mock)
}

func (poker_table_db POKER_TABLE_DB) TestProceedNextRound() {
	mock := Models.MOCK_REQUEST{
		Table_name:     "test",
		Table_passcode: "test",
	}
	poker_table_db.ProcessProceedToNextRoundRequest(mock)
}

// END TESTING
