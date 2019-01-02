package entities

import "fmt"

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

// String represents status
func (a *Agent) String() string {
	return fmt.Sprintf("%s(%d):%v:%v", Meta.Dynamic[AGENT].Short,
		a.Human.ID, a.Current, a.Human.Pos())
}
