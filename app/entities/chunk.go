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

func (ch *Chunk) Init(m *Model) *Chunk {
	ch.Base.Init(CHUNK, m)
	ch.Shape.P1 = &ch.Point
	ch.InRailEdges = make(map[uint]*DelegateRailEdge)
	ch.OutRailEdges = make(map[uint]*DelegateRailEdge)
	return ch
}

func (ch *Chunk) Add(raw Entity) {
	switch obj := raw.(type) {
	case *Track:
		ch.addTrack(obj)
	}
}

func (ch *Chunk) addTrack(tr *Track) {
	if ch.RailNode == nil {
		ch.RailNode = &DelegateRailNode{
			Base:  ch.M.NewBase(RAILNODE, tr.O),
			Point: tr.FromNode.Point,
		}
		ch.RailNode.RailNodes = make(map[uint]*RailNode)
	}
	ch.RailNode.RailNodes[tr.FromNode.ID] = tr.FromNode

	target := ch.M.RootCluster.FindChunk(tr.ToNode, ch.Parent.Scale)

	if ch.OutRailEdges[target.ID] == nil {
		re := &DelegateRailEdge{
			Base:   ch.M.NewBase(RAILEDGE, tr.O),
			From:   ch.RailNode,
			FromID: ch.RailNode.ID,
			To:     target.RailNode,
			ToID:   target.RailNode.ID,
			Tracks: make(map[uint]*Track),
		}
		re.Tracks[tr.ID] = tr
		ch.OutRailEdges[target.ID] = re
		target.InRailEdges[ch.ID] = re
		if reverse, ok := target.OutRailEdges[ch.ID]; ok {
			re.ReverseID = reverse.ID
			reverse.ReverseID = re.ID
		}
	}
}

func (ch *Chunk) Remove(raw Entity) {
	switch obj := raw.(type) {
	case *Track:
		for _, re := range ch.OutRailEdges {
			delete(re.Tracks, obj.ID)
			if len(re.Tracks) == 0 {
				delete(ch.OutRailEdges, re.ID)
			}
		}
		delete(ch.RailNode.RailNodes, obj.FromNode.ID)
		if len(ch.RailNode.RailNodes) == 0 {
			ch.RailNode = nil
		}
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
