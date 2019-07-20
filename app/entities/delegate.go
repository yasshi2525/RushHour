package entities

import "encoding/json"

type DelegateMap struct {
	RailNodes jsonDelegateRailNode `json:"rail_nodes"`
	RailEdges jsonDelegateRailEdge `json:"rail_edges"`
}

func (dm *DelegateMap) Init() {
	dm.RailNodes = make(map[uint]*DelegateRailNode)
	dm.RailEdges = make(map[uint]*DelegateRailEdge)
}

type DelegateRailNode struct {
	Base
	Point

	RailNodes map[uint]*RailNode `json:"-"`
}

type jsonDelegateRailNode map[uint]*DelegateRailNode

func (jrn jsonDelegateRailNode) MarshalJSON() ([]byte, error) {
	rns := make([]*DelegateRailNode, len(jrn))
	for i, rn := range jrn {
		rns[i] = rn
	}
	return json.Marshal(rns)
}

type DelegateRailEdge struct {
	Base

	From *DelegateRailNode `json:"-"`
	To   *DelegateRailNode `json:"-"`

	Tracks map[uint]*Track `json:"-"`

	FromID    uint `json:"from"`
	ToID      uint `json:"to"`
	ReverseID uint `json:"eid"`
}

type jsonDelegateRailEdge map[uint]*DelegateRailEdge

func (jre jsonDelegateRailEdge) MarshalJSON() ([]byte, error) {
	res := make([]*DelegateRailEdge, len(jre))
	for i, re := range jre {
		res[i] = re
	}
	return json.Marshal(res)
}
