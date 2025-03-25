package table

type CardLayout string

const (
	Sequential = "seq"
	Fibonacci  = "fib"
)

type Card struct {
	Value int
}

type Results_Card struct {
	Value     int
	Winner    bool
	VoteCount int
}
