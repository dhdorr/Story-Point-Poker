package Models

type ROUNDS_DB struct {
	rounds []ROUND
}

type ROUND struct {
	votes map[string]VOTE
}
