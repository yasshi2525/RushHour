package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DelegateMap represents Map for client view.
type DelegateMap struct {
	RailNodes jsonDelegateRailNode `json:"rail_nodes"`
	RailEdges jsonDelegateRailEdge `json:"rail_edges"`
}

// Init creates maps.
func (dm *DelegateMap) Init(m *Model) {
	dm.RailNodes = make(map[uint]*DelegateRailNode)
	dm.RailEdges = make(map[uint]*DelegateRailEdge)
}

// JSONPlayer is collection of Player.
type JSONPlayer map[uint]*Player

// MarshalJSON serializes collection of Player.
func (jp JSONPlayer) MarshalJSON() ([]byte, error) {
	os := []*Player{}
	for _, o := range jp {
		os = append(os, o)
	}
	return json.Marshal(os)
}

// DelegateRailNode is delegate RailNode.
type DelegateRailNode struct {
	Base

	RailNodes map[uint]*RailNode `json:"-"`
	Pos       *Point             `json:"pos"`
	Multi     int                `json:"mul"`
	ParentID  uint               `json:"pid,omitempty"`
	ChildID   uint               `json:"cid,omitempty"`
}

// UpdatePos update center point.
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

// DelegateRailEdge is delegate of RailEdge
type DelegateRailEdge struct {
	Base

	From    *DelegateRailNode `json:"-"`
	To      *DelegateRailNode `json:"-"`
	Reverse *DelegateRailEdge `json:"-"`

	RailEdges map[uint]*RailEdge `json:"-"`
	Multi     int                `json:"mul"`

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
