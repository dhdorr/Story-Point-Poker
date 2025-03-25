package card

type Card struct {
	Value int
}

type Vote struct {
	PlayerID string
	Card     Card
}
