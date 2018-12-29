package entities

// Agent is the wrapper of Human.
// Agent can move concerting minimum cost route.
type Agent struct {
	Human
	Current Edge
}

func (a *Agent) Consume(interval float64) {
	//TODO
}
