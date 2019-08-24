package entities

import (
	"fmt"
)

type Chunk struct {
	Base
	Shape
	Point

	RailNode *DelegateRailNode

	Parent *Cluster

	InRailEdges  map[uint]*DelegateRailEdge
	OutRailEdges map[uint]*DelegateRailEdge
}

func (m *Model) NewChunk(p *Cluster, o *Player) *Chunk {
	ch := &Chunk{
		Base:  m.NewBase(CHUNK, o),
		Point: p.Point,
	}
	ch.Init(m)
	ch.Resolve(p)
	m.Add(ch)
	return ch
}

func (ch *Chunk) B() *Base {
	return &ch.Base
}

func (ch *Chunk) S() *Shape {
	return &ch.Shape
}

func (ch *Chunk) Init(m *Model) {
	ch.Base.Init(CHUNK, m)
	ch.Shape.P1 = &ch.Point
	ch.InRailEdges = make(map[uint]*DelegateRailEdge)
	ch.OutRailEdges = make(map[uint]*DelegateRailEdge)
}

func (ch *Chunk) Add(raw Entity) {
	switch obj := raw.(type) {
	case *RailNode:
		ch.addRailNode(obj)
	case *RailEdge:
		ch.addRailEdge(obj)
	}
}

func (ch *Chunk) addRailNode(rn *RailNode) {
	if ch.RailNode == nil {
		var pid uint
		if parent := ch.Parent.Parent; parent != nil && parent.Data[rn.OwnerID] != nil {
			p := parent.Data[rn.OwnerID].RailNode
			pid = p.ID
		}
		ch.RailNode = &DelegateRailNode{
			Base:     ch.M.NewBase(RAILNODE, rn.O),
			Pos:      &rn.Point,
			ParentID: pid,
		}
		ch.RailNode.RailNodes = make(map[uint]*RailNode)
	}
	ch.RailNode.RailNodes[rn.ID] = rn
	ch.RailNode.Multi = len(ch.RailNode.RailNodes)
	if ch.RailNode.Multi == 1 {
		ch.RailNode.ChildID = rn.ID
	} else {
		ch.RailNode.ChildID = 0
	}
	ch.RailNode.UpdatePos()
}

func (ch *Chunk) addRailEdge(re *RailEdge) {
	target := ch.M.RootCluster.FindChunk(re.ToNode, ch.Parent.Scale)

	if ch.OutRailEdges[target.ID] == nil {
		dre := &DelegateRailEdge{
			Base:      ch.M.NewBase(RAILEDGE, re.O),
			From:      ch.RailNode,
			FromID:    ch.RailNode.ID,
			To:        target.RailNode,
			ToID:      target.RailNode.ID,
			RailEdges: make(map[uint]*RailEdge),
		}
		ch.OutRailEdges[target.ID] = dre
		target.InRailEdges[ch.ID] = dre
		if reverse, ok := target.OutRailEdges[ch.ID]; ok {
			dre.ReverseID = reverse.ID
			reverse.ReverseID = dre.ID
		}
	}
	dre := ch.OutRailEdges[target.ID]
	dre.RailEdges[re.ID] = re
	dre.Multi = len(dre.RailEdges)
}

func (ch *Chunk) Remove(raw Entity) {
	switch obj := raw.(type) {
	case *RailNode:
		ch.removeRailNode(obj)
	case *RailEdge:
		ch.removeRailEdge(obj)
	}
}

func (ch *Chunk) removeRailNode(rn *RailNode) {
	delete(ch.RailNode.RailNodes, rn.ID)
	ch.RailNode.UpdatePos()
	ch.RailNode.Multi = len(ch.RailNode.RailNodes)

	if ch.RailNode.Multi == 1 {
		ch.RailNode.ChildID = rn.ID
	} else {
		ch.RailNode.ChildID = 0
	}

	if len(ch.RailNode.RailNodes) == 0 {
		ch.RailNode = nil
	}
}

func (ch *Chunk) removeRailEdge(re *RailEdge) {
	target := ch.M.RootCluster.FindChunk(re.ToNode, ch.Parent.Scale)
	dre := ch.OutRailEdges[target.ID]
	delete(dre.RailEdges, re.ID)
	dre.Multi = len(dre.RailEdges)
	if len(dre.RailEdges) == 0 {
		delete(ch.OutRailEdges, dre.ID)
		delete(target.InRailEdges, dre.ID)
	}
}

func (ch *Chunk) Has(raw Entity) bool {
	switch obj := raw.(type) {
	case *RailNode:
		if ch.RailNode == nil {
			return false
		}
		_, ok := ch.RailNode.RailNodes[obj.ID]
		return ok
	case *RailEdge:
		_, ok := ch.OutRailEdges[obj.ID]
		return ok
	}
	return false
}

func (ch *Chunk) IsEmpty() bool {
	return ch.RailNode == nil
}

func (ch *Chunk) CheckDelete() error {
	return nil
}

func (ch *Chunk) BeforeDelete() {
	ch.Parent.UnResolve(ch)
}

func (ch *Chunk) Delete() {
	ch.M.Delete(ch)
}

func (ch *Chunk) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Cluster:
			ch.Parent = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

func (ch *Chunk) Export(dm *DelegateMap) {
	if rn := ch.RailNode; rn != nil {
		dm.RailNodes[rn.ID] = rn
	}
	for _, re := range ch.InRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.FromID] = re.From
	}
	for _, re := range ch.OutRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.ToID] = re.To
	}
}

// String represents status
func (ch *Chunk) String() string {
	return fmt.Sprintf("%s(%.1f:%d):u=%d,%v,i=%d,o=%d:%v", ch.Type().Short(),
		ch.Parent.Scale, ch.ID, ch.OwnerID, ch.RailNode, len(ch.InRailEdges), len(ch.OutRailEdges), ch.Point)
}
