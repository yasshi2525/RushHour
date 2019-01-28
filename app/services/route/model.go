package route

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// Model has minimux distance route information to specific Node.
type Model struct {
	GoalIDs []uint
	Nodes   map[entities.ModelType]map[uint]*Node
	Edges   map[entities.ModelType]map[uint]*Edge
}

// NewModel creates instance or copies original if it is specified.
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

// Export copies Nodes and Edges having same id.
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

// ExportWith copies Nodes and Edges, then returns corresponding Node to specified goal.
func (m *Model) ExportWith(t entities.ModelType, id uint) (*Model, *Node) {
	copy := m.Export()
	return copy, copy.Nodes[t][id]
}

// NumNodes returns the number of Nodes.
func (m *Model) NumNodes() int {
	var sum int
	for _, ns := range m.Nodes {
		sum += len(ns)
	}
	return sum
}

// NumEdges returns the number of Edges.
func (m *Model) NumEdges() int {
	var sum int
	for _, es := range m.Edges {
		sum += len(es)
	}
	return sum
}

// AddGoalID adds id as goal.
func (m *Model) AddGoalID(id uint) {
	m.GoalIDs = append(m.GoalIDs, id)
}

// FindOrCreateNode returns corresponding Node.
// If such Node doesn't exist, create and return new Node.
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

// FindOrCreateEdge returns corresponding Edge.
// If such Edge doesn't exist, create and return new Edge.
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

// Fix discards no more using data.
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

// Payload is collection of Model.
type Payload struct {
	Route     map[uint]*Model
	Processed int
	Total     int
}

// IsOK returns whether all Model was built or not.
func (p *Payload) IsOK() bool {
	return p.Processed == p.Total
}

// Import accepts result of calcuration.
func (p *Payload) Import(oth *Payload) {
	for goalID, model := range oth.Route {
		p.Route[goalID] = model
	}
	p.Processed++
}
