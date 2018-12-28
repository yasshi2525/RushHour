package routers

import "github.com/yasshi2525/RushHour/app/models/entities"

// Agent is the wrapper of Human.
// Agent can move concerting minimum cost route.
type Agent struct {
	entities.Human
	Current Edge
}

func (a *Agent) Consume(interval float64) {
	//TODO
}
