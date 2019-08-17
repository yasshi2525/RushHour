package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DelegateMap struct {
	Players   jsonPlayer           `json:"players"`
	RailNodes jsonDelegateRailNode `json:"rail_nodes"`
	RailEdges jsonDelegateRailEdge `json:"rail_edges"`
}

func (dm *DelegateMap) Init(m *Model) {
	dm.Players = m.Players
	dm.RailNodes = make(map[uint]*DelegateRailNode)
	dm.RailEdges = make(map[uint]*DelegateRailEdge)
}

type jsonPlayer map[uint]*Player

func (jp jsonPlayer) MarshalJSON() ([]byte, error) {
	os := []*Player{}
	for _, o := range jp {
		os = append(os, o)
	}
	return json.Marshal(os)
}

type DelegateRailNode struct {
	Base

	RailNodes map[uint]*RailNode `json:"-"`
	Pos       *Point             `json:"pos"`
	Multi     int                `json:"mul"`
	Color     int                `json:"color"`
	ParentID  uint               `json:"pid"`
}

func (rn *DelegateRailNode) UpdatePos() {
	rn.Pos = &Point{}

	for _, child := range rn.RailNodes {
		rn.Pos.X += child.X / float64(len(rn.RailNodes))
		rn.Pos.Y += child.Y / float64(len(rn.RailNodes))
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
	Color     int                `json:"color"`

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
