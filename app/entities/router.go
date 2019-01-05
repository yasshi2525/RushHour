package entities

import (
	"math"
	"sort"
)

// Digest directs resource of model
type Digest struct {
	Type ModelType
	ID   uint
}

// Node is digest of Relayable, Transportable for routing.
// The chain of Node represents one route.
type Node struct {
	Base    *Digest
	Cost    float64
	Via     *Node
	ViaEdge *Edge
	Out     []*Edge
	In      []*Edge
}

// NewNode returns instance
func NewNode(t ModelType, id uint) *Node {
	return &Node{
		Base: &Digest{t, id},
		Cost: math.MaxFloat64,
		Out:  []*Edge{},
		In:   []*Edge{},
	}
}

// SameAs check both directs same resource
func (n *Node) SameAs(oth Indexable) bool {
	return n.Base.Type == oth.Type() && n.Base.ID == oth.Idx()
}

// WalkThrough set distance towrards self to Cost of connected Nodes.
// Initial cost of connected Node must be max float64 value.
func (n *Node) WalkThrough() {
	var x *Node
	var q nodeQueue = []*Node{n}
	n.Cost = 0

	for len(q) > 0 {
		x, q = q[0], q[1:]

		for _, e := range x.In {
			y := e.From
			v := x.Cost + e.Cost
			if v < y.Cost {
				y.Cost = v
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
		if e.To == n.Via {
			n.ViaEdge = e
			break
		}
	}
	// in order to save memory
	n.In = nil
	n.Out = nil
}

// Edge is wrapper of Step for routing.
type Edge struct {
	Base *Digest
	From *Node
	To   *Node
	Cost float64
}

// NewEdge creates instance and append slice of Node
func NewEdge(base *Digest, from *Node, to *Node, v float64) *Edge {
	e := &Edge{
		Base: base,
		From: from,
		To:   to,
		Cost: v,
	}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
	return e
}

// GenStepEdges generates Edge list from Nodes and Steps.
func GenStepEdges(ns []*Node, steps map[uint]*Step) []*Edge {
	es := []*Edge{}
	for _, s := range steps {
		es = append(es, genEdge(ns, s))
	}
	return es
}

// GenTrackEdges generates Edge list from Nodes and LineTask.
func GenLineTaskEdges(ns []*Node, lts map[uint]*LineTask) []*Edge {
	es := []*Edge{}
	for _, lt := range lts {
		es = append(es, genEdge(ns, lt))
	}
	return es
}

func genEdge(ns []*Node, base Connectable) *Edge {
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
	return NewEdge(&Digest{base.Type(), base.Idx()}, from, to, base.Cost())
}

// nodeQueue is open list for searching
type nodeQueue []*Node

func (q nodeQueue) Len() int {
	return len(q)
}

func (q nodeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q nodeQueue) Less(i, j int) bool {
	return q[i].Cost < q[j].Cost
}
