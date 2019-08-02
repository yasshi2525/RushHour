package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
	Multi     int                `json:"mul"`
	Scale     float64            `json:"sc"`
}

func (rn *DelegateRailNode) UpdatePos() {
	rn.X, rn.Y = 0, 0

	for _, child := range rn.RailNodes {
		rn.X += child.X / float64(len(rn.RailNodes))
		rn.Y += child.Y / float64(len(rn.RailNodes))
	}
}

// String represents status
func (rn *DelegateRailNode) String() string {
	list := []string{}

	for id := range rn.RailNodes {
		list = append(list, fmt.Sprintf("rn(%d)", id))
	}

	return fmt.Sprintf("^rn(%d):%v", rn.ID,
		strings.Join(list, ","))
}

type jsonDelegateRailNode map[uint]*DelegateRailNode

func (jrn jsonDelegateRailNode) MarshalJSON() ([]byte, error) {
	rns := []*DelegateRailNode{}
	for _, rn := range jrn {
		rns = append(rns, rn)
	}
	return json.Marshal(rns)
}

type DelegateRailEdge struct {
	Base

	From *DelegateRailNode `json:"-"`
	To   *DelegateRailNode `json:"-"`

	RailEdges map[uint]*RailEdge `json:"-"`
	Multi     int                `json:"mul"`
	Scale     float64            `json:"sc"`

	FromID    uint `json:"from"`
	ToID      uint `json:"to"`
	ReverseID uint `json:"eid"`
}

type jsonDelegateRailEdge map[uint]*DelegateRailEdge

func (jre jsonDelegateRailEdge) MarshalJSON() ([]byte, error) {
	res := []*DelegateRailEdge{}
	for _, re := range jre {
		res = append(res, re)
	}
	return json.Marshal(res)
}
