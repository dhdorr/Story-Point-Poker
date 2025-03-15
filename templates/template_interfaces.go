package templates

import "dhdorr/story-point-poker/table"

type Gen_Test_A struct {
	Data int
}

type Gen_Test_B struct {
	Data string
}

type Gen_Test_Interface interface {
	Gen_Test_A | Gen_Test_B
}

type Gen_Table_Session_Interface interface {
	table.Table_Session | table.Table_Session_Constructor | Gen_Test_B
}
