package route

import (
	"fmt"
	"math"
	"sort"

	"github.com/yasshi2525/RushHour/app/entities"
)

// Digest directs resource of model
type Digest struct {
	ModelType entities.ModelType
	ID        uint
	Value     float64
}

// NodeQueue is open list for searching
type NodeQueue []*Node

func (q NodeQueue) Len() int {
	return len(q)
}

func (q NodeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q NodeQueue) Less(i, j int) bool {
	return q[i].Value < q[j].Value
}

// Node is digest of Relayable, Transportable for routing.
// The chain of Node represents one route.
type Node struct {
	Digest
	Via     *Node
	ViaEdge *Edge
	Out     []*Edge
	In      []*Edge
}

// NewNode returns instance
func NewNode(obj entities.Indexable) *Node {
	return &Node{
		Digest: Digest{obj.Type(), obj.Idx(), math.MaxFloat64},
		Out:    []*Edge{},
		In:     []*Edge{},
	}
}

func (n *Node) Idx() uint {
	return n.ID
}

func (n *Node) Type() entities.ModelType {
	return n.ModelType
}

// SameAs check both directs same resource
func (n *Node) SameAs(oth entities.Indexable) bool {
	return n.ModelType == oth.Type() && n.ID == oth.Idx()
}

// WalkThrough set distance towrards self to Value of connected Nodes.
// Initial cost of connected Node must be max float64 value.
func (n *Node) WalkThrough() {
	var x *Node
	var q NodeQueue = []*Node{n}
	n.Value = 0

	for len(q) > 0 {
		x, q = q[0], q[1:]

		for _, e := range x.In {
			y := e.FromNode
			v := x.Value + e.Cost()
			if v < y.Value {
				y.Value = v
				y.Via = x
				q = append(q, y)
				sort.Sort(q)
			}
		}
	}
}

// Fix sets ViaEdge and discards no more need slice
func (n *Node) Fix() {
	for _, e := range n.Out {
		if e.ToNode == n.Via {
			n.ViaEdge = e
			break
		}
	}
	// in order to save memory
	n.In = nil
	n.Out = nil
}

func (n *Node) String() string {
	viastr := ""
	if n.ViaEdge != nil {
		viastr = fmt.Sprintf(",via=%d(->%d)", n.ViaEdge.Idx(), n.ViaEdge.To().Idx())
	}
	valstr := ""
	if n.Value == math.MaxFloat64 {
		valstr = "NaN"
	} else {
		valstr = fmt.Sprintf("%.2f", n.Value)
	}
	return fmt.Sprintf("Node(%v,%d):i=%d,o=%d,v=%s%s",
		n.Type(), n.Idx(), len(n.In), len(n.Out), valstr, viastr)
}
