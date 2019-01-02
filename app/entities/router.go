package entities

import (
	"fmt"
	"math"
	"sort"
)

// Node is wrapper of Relayable for routing.
// The chain of Node represents one route.
type Node struct {
	Base    Relayable
	Cost    float64
	Via     *Node
	ViaEdge *Edge
	Out     []*Edge
	In      []*Edge
}

// NewNode returns instance
func NewNode(base Relayable) *Node {
	return &Node{
		Base: base,
		Cost: math.MaxFloat64,
		Out:  []*Edge{},
		In:   []*Edge{},
	}
}

// Edge is wrapper of Step for routing.
type Edge struct {
	Base *Step
	From *Node
	To   *Node
}

// NewEdge creates instance and append slice of Node
func NewEdge(base *Step, from *Node, to *Node) *Edge {
	e := &Edge{
		Base: base,
		From: from,
		To:   to,
	}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
	return e
}

// GenEdges generates Edge list from Nodes and Steps.
func GenEdges(ns []*Node, steps map[uint]*Step) []*Edge {
	es := []*Edge{}
	for _, s := range steps {
		var from, to *Node
		for _, n := range ns {
			if n.Base == s.From() {
				from = n
			}
			if n.Base == s.To() {
				to = n
			}
			if from != nil && to != nil {
				break
			}
		}
		if from == nil && to == nil {
			panic(fmt.Errorf("fail to create edge from %v: from=%v, to=%v", s, from, to))
		}
		es = append(es, NewEdge(s, from, to))
	}
	return es
}

// Cost is evaluated for minium cost searching.
func (e *Edge) Cost() float64 {
	return e.Base.Cost()
}

type NodeQueue []*Node

func (q NodeQueue) Len() int {
	return len(q)
}

func (q NodeQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q NodeQueue) Less(i, j int) bool {
	return q[i].Cost < q[j].Cost
}

// WalkThrough set distance towrards self to Cost of connected Nodes.
// Initial cost of connected Node must be max float64 value.
func (n *Node) WalkThrough() {
	var x *Node
	var q NodeQueue = []*Node{n}
	n.Cost = 0

	for len(q) > 0 {
		x, q = q[0], q[1:]

		for _, e := range x.In {
			y := e.From
			v := x.Cost + e.Cost()
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
