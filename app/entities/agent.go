package entities

// Agent is the wrapper of Human.
// Agent can move concerting minimum cost route.
type Agent struct {
	Human   *Human
	Current *Edge
}

// NewAgent creates instance
func NewAgent(h *Human) *Agent {
	return &Agent{
		Human: h,
	}
}

// Consume makes Human moves
func (a *Agent) Consume(interval float64) {
	//TODO
}
