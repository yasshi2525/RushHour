package route

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

type Model struct {
	GoalIDs []uint
	Nodes   map[entities.ModelType]map[uint]*Node
	Edges   map[entities.ModelType]map[uint]*Edge
}

func NewModel(origin ...*Model) *Model {
	nodes := make(map[entities.ModelType]map[uint]*Node)
	edges := make(map[entities.ModelType]map[uint]*Edge)

	if len(origin) > 0 { //copy
		origin := origin[0]
		goalIDs := make([]uint, len(origin.GoalIDs))
		for i, goalID := range origin.GoalIDs {
			goalIDs[i] = goalID
		}
		for key := range origin.Nodes {
			nodes[key] = make(map[uint]*Node)
		}
		for key := range origin.Edges {
			edges[key] = make(map[uint]*Edge)
		}
		return &Model{goalIDs, nodes, edges}
	}
	return &Model{[]uint{}, nodes, edges}
}

func (m *Model) Export() *Model {
	copy := NewModel(m)
	for res, ns := range m.Nodes {
		for id, n := range ns {
			copy.Nodes[res][id] = n.Export()
		}
	}
	for res, es := range m.Edges {
		for id, e := range es {
			oldFrom, oldTo := e.FromNode, e.ToNode
			newFrom := copy.Nodes[oldFrom.ModelType][oldFrom.ID]
			newTo := copy.Nodes[oldTo.ModelType][oldTo.ID]
			copy.Edges[res][id] = e.Export(newFrom, newTo)
		}
	}
	return copy
}

func (m *Model) ExportWith(t entities.ModelType, id uint) (*Model, *Node) {
	copy := m.Export()
	return copy, copy.Nodes[t][id]
}

func (m *Model) NumNodes() int {
	var sum int
	for _, ns := range m.Nodes {
		sum += len(ns)
	}
	return sum
}

func (m *Model) NumEdges() int {
	var sum int
	for _, es := range m.Edges {
		sum += len(es)
	}
	return sum
}

func (m *Model) AddGoalID(id uint) {
	m.GoalIDs = append(m.GoalIDs, id)
}

func (m *Model) FindOrCreateNode(origin entities.Entity) *Node {
	if _, ok := m.Nodes[origin.B().Type()]; !ok {
		m.Nodes[origin.B().Type()] = make(map[uint]*Node)
	}
	if n, ok := m.Nodes[origin.B().Type()][origin.B().Idx()]; ok {
		return n
	}
	n := NewNode(origin)
	m.Nodes[origin.B().Type()][origin.B().Idx()] = n
	return n
}

func (m *Model) FindOrCreateEdge(origin entities.Connectable) *Edge {
	if _, ok := m.Edges[origin.B().Type()]; !ok {
		m.Edges[origin.B().Type()] = make(map[uint]*Edge)
	}
	if e, ok := m.Edges[origin.B().Type()][origin.B().Idx()]; ok {
		return e
	}
	from, to := m.FindOrCreateNode(origin.From()), m.FindOrCreateNode(origin.To())
	e := NewEdge(origin, from, to)
	m.Edges[origin.B().Type()][origin.B().Idx()] = e
	return e
}

func (m *Model) Fix() {
	for _, ns := range m.Nodes {
		for _, n := range ns {
			n.Fix()
		}
	}
}

func (m *Model) String() string {
	nodes := make(map[entities.ModelType][]uint)
	for res, ns := range m.Nodes {
		nodes[res] = make([]uint, len(ns))
		var i int
		for id := range ns {
			nodes[res][i] = id
			i++
		}
	}
	edges := make(map[entities.ModelType][]uint)
	for res, es := range m.Edges {
		edges[res] = make([]uint, len(es))
		var i int
		for id := range es {
			edges[res][i] = id
			i++
		}
	}
	return fmt.Sprintf("Route:g=%v,n=%v,e=%v", m.GoalIDs, nodes, edges)
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
	for goalID, model := range oth.Route {
		p.Route[goalID] = model
	}
	p.Processed++
}
