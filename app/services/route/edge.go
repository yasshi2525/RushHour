package route

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// Edge is wrapper of Step for routing.
type Edge struct {
	Digest
	FromNode *Node
	ToNode   *Node
}

func NewEdge(base entities.Connectable, from *Node, to *Node) *Edge {
	e := &Edge{Digest{base.Type(), base.Idx(), base.Cost()}, from, to}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
	return e
}

// NewEdgeFrom creates instance and append slice of Node
func NewEdgeFrom(ns []*Node, base entities.Connectable) *Edge {
	var from, to *Node
	for _, n := range ns {
		if n.SameAs(base.From()) {
			from = n
		}
		if n.SameAs(base.To()) {
			to = n
		}
		if from != nil && to != nil {
			break
		}
	}
	return NewEdge(base, from, to)
}

func (e *Edge) Idx() uint {
	return e.ID
}

func (e *Edge) Type() entities.ModelType {
	return e.ModelType
}

func (e *Edge) From() entities.Indexable {
	return e.FromNode
}

func (e *Edge) To() entities.Indexable {
	return e.ToNode
}

func (e *Edge) Cost() float64 {
	return e.Value
}

func (e *Edge) Export(ns []*Node) *Edge {
	var from, to *Node
	for _, n := range ns {
		if n.SameAs(e.From()) {
			from = n
		}
		if n.SameAs(e.To()) {
			to = n
		}
		if from != nil && to != nil {
			break
		}
	}
	return NewEdge(e, from, to)
}

func (e *Edge) String() string {
	return fmt.Sprintf("Edge(%v,%d):f=%d,t=%d,v=%.2f",
		e.Type(), e.Idx(), e.From().Idx(), e.To().Idx(), e.Cost())
}
