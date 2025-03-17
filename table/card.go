package table

type CardLayout string

const (
	Sequential = "seq"
	Fibonacci  = "fib"
)

type Card struct {
	Value int
}
