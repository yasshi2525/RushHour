package entities

// Step represents two Junction is logically connected.
// Step is out of target for persistence because it can intruduce by other resources.
type Step struct {
	ID     uint
	From   *Junction
	To     *Junction
	Weight float64
}

// Cost is calculated by distance * weight of Step
func (s *Step) Cost() float64 {
	return s.From.Dist(&s.To.Point) * s.Weight
}

// Node is wrapper of Junction for routing.
// The chain of Node represents one route.
type Node struct {
	Original *Junction
	Cost     float64
	Via      *Node
	Out      []*Edge
	In       []*Edge
}

// Edge is wrapper of Step for routing.
type Edge struct {
	Original *Step
	From     *Node
	To       *Node
}

// Cost is evaluated for minium cost searching.
func (e *Edge) Cost() float64 {
	return e.Original.Cost()
}
