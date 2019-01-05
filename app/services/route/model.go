package route

import (
	"github.com/yasshi2525/RushHour/app/entities"
)

type Model struct {
	GoalIDs []uint
	Nodes   []*Node
	Edges   []*Edge
}

func NewModel(origin *Model) *Model {
	cids := make([]uint, len(origin.GoalIDs))
	for i, cid := range origin.GoalIDs {
		cids[i] = cid
	}
	ns := make([]*Node, len(origin.Nodes))
	es := make([]*Edge, len(origin.Edges))
	return &Model{cids, ns, es}
}

func (m *Model) Export() *Model {
	copy := NewModel(m)
	for i, n := range m.Nodes {
		copy.Nodes[i] = NewNode(n)
	}
	for i, e := range m.Edges {
		copy.Edges[i] = e.Export(copy.Nodes)
	}
	return copy
}

func (m *Model) ExportWith(t entities.ModelType, id uint) (*Model, *Node) {
	copy := NewModel(m)
	var p *Node
	for i, n := range m.Nodes {
		newN := NewNode(n)
		if newN.Type() == t && newN.Idx() == id {
			p = newN
		}
		copy.Nodes[i] = newN
	}
	for i, e := range m.Edges {
		copy.Edges[i] = e.Export(copy.Nodes)
	}
	return copy, p
}

type Payload struct {
	Route     map[uint]*Model
	Processed int
	Total     int
}

func (p *Payload) IsOK() bool {
	return p.Processed == p.Total
}

func (p *Payload) Import(oth *Payload) {
	for cid, model := range oth.Route {
		p.Route[cid] = model
	}
	p.Processed++
}
