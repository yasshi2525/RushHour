package routers

import "github.com/yasshi2525/RushHour/app/models/entities"

// Node is wrapper of Junction for routing.
// The chain of Node represents one route.
type Node struct {
	Original *entities.Junction
	Cost     float64
	Via      *Node
	Out      []*Edge
	In       []*Edge
}

// Edge is wrapper of Step for routing.
type Edge struct {
	Original *entities.Step
	From     *Node
	To       *Node
}

// Cost is evaluated for minium cost searching.
func (e *Edge) Cost() float64 {
	return e.Original.Cost()
}
