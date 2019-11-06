package route

import (
	"fmt"

	"github.com/yasshi2525/RushHour/entities"
)

// Edge is wrapper of Step for routing.
type Edge struct {
	Digest
	FromNode *Node
	ToNode   *Node
}

// NewEdge creates new instance.
func NewEdge(base entities.Connectable, from *Node, to *Node) *Edge {
	e := &Edge{Digest{base.B().Type(), base.B().Idx(), base.Cost()}, from, to}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
	return e
}

// Export creates new instance which has specified Node.
func (e *Edge) Export(from *Node, to *Node) *Edge {
	newE := &Edge{e.Digest, from, to}
	from.Out = append(from.Out, newE)
	to.In = append(to.In, newE)
	return newE
}

// Cost returns value.
func (e *Edge) Cost() float64 {
	return e.Value
}

func (e *Edge) String() string {
	return fmt.Sprintf("Edge(%v,%d):f=%d,t=%d,v=%.2f",
		e.ModelType, e.ID, e.FromNode.ID, e.ToNode.ID, e.Cost())
}
