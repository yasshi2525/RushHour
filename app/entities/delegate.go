package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DelegateMap represents Map for client view.
type DelegateMap struct {
	Residences jsonDelegateResidence `json:"residences"`
	Compannies jsonDelegateCompany   `json:"companies"`
	RailNodes  jsonDelegateRailNode  `json:"rail_nodes"`
	RailEdges  jsonDelegateRailEdge  `json:"rail_edges"`
}

// Init creates maps.
func (dm *DelegateMap) Init(m *Model) {
	dm.Residences = make(map[uint]*DelegateResidence)
	dm.Compannies = make(map[uint]*DelegateCompany)
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

type delegateLocalable interface {
	B() *Base
}

// DelegateNode represents delegation Localable.
type DelegateNode struct {
	Base

	Pos      *Point             `json:"pos"`
	Multi    int                `json:"mul"`
	ParentID uint               `json:"pid,omitempty"`
	ChildID  uint               `json:"cid,omitempty"`
	List     map[uint]Localable `json:"-"`
}

// NewDelegateNode creates instance.
func (ch *Chunk) NewDelegateNode(obj Localable, pid uint) DelegateNode {
	return DelegateNode{
		Base:     ch.M.NewBase(obj.B().T, obj.B().O),
		Pos:      obj.Pos().Clone(),
		ParentID: pid,
		List:     make(map[uint]Localable),
	}
}

// Add accepts new instance ant increment count variable
func (dn *DelegateNode) Add(obj Localable) {
	dn.List[obj.B().ID] = obj
	dn.updateMulti()
	dn.updatePos()
}

// Remove delete argument from list ant decrement count variable
func (dn *DelegateNode) Remove(obj Localable) {
	delete(dn.List, obj.B().ID)
	dn.updateMulti()
	dn.updatePos()
}

func (dn *DelegateNode) updateMulti() {
	dn.Multi = len(dn.List)
	if dn.Multi == 1 {
		for idx := range dn.List {
			dn.ChildID = idx
		}
	} else {
		dn.ChildID = 0
	}
}

func (dn *DelegateNode) updatePos() {
	dn.Pos = &Point{}

	for _, child := range dn.List {
		dn.Pos.X += child.Pos().X / float64(len(dn.List))
		dn.Pos.Y += child.Pos().Y / float64(len(dn.List))
	}
}

// String represents status
func (dn *DelegateNode) String() string {
	list := []string{}

	for id := range dn.List {
		list = append(list, fmt.Sprintf("%s(%d)", dn.T.Short(), id))
	}

	return fmt.Sprintf("^%s(%d):%v", dn.T.Short(), dn.ID, strings.Join(list, ","))
}

// DelegateEdge represents delegation Connectable.
type DelegateEdge struct {
	Base

	Multi int                  `json:"mul"`
	List  map[uint]Connectable `json:"-"`
	From  delegateLocalable    `json:"-"`
	To    delegateLocalable    `json:"-"`

	FromID  uint `json:"from"`
	ToID    uint `json:"to"`
	ChildID uint `json:"cid,omitempty"`
}

// NewDelegateEdge creates instance.
func (ch *Chunk) NewDelegateEdge(obj Connectable, from delegateLocalable, to delegateLocalable) DelegateEdge {
	return DelegateEdge{
		Base:   ch.M.NewBase(obj.B().T, obj.B().O),
		List:   make(map[uint]Connectable),
		From:   from,
		FromID: from.B().ID,
		To:     to,
		ToID:   to.B().ID,
	}
}

// Add accepts new instance ant increment count variable
func (de *DelegateEdge) Add(obj Connectable) {
	de.List[obj.B().ID] = obj
	de.updateMulti()
}

// Remove delete argument from list ant decrement count variable
func (de *DelegateEdge) Remove(obj Connectable) {
	delete(de.List, obj.B().ID)
	de.updateMulti()
}

func (de *DelegateEdge) updateMulti() {
	de.Multi = len(de.List)
	if de.Multi == 1 {
		for idx := range de.List {
			de.ChildID = idx
		}
	} else {
		de.ChildID = 0
	}
}

// DelegateResidence is delegate of Residence
type DelegateResidence struct {
	DelegateNode
}

// B returns reference of Base Object
func (dr DelegateResidence) B() *Base {
	return &dr.Base
}

type jsonDelegateResidence map[uint]*DelegateResidence

func (jr jsonDelegateResidence) MarshalJSON() ([]byte, error) {
	rs := []*DelegateResidence{}
	for _, r := range jr {
		rs = append(rs, r)
	}
	return json.Marshal(rs)
}

// DelegateCompany is delegate of RailNode
type DelegateCompany struct {
	DelegateNode
}

// B returns reference of Base Object
func (dc DelegateCompany) B() *Base {
	return &dc.Base
}

type jsonDelegateCompany map[uint]*DelegateCompany

func (jc jsonDelegateCompany) MarshalJSON() ([]byte, error) {
	cs := []*DelegateCompany{}
	for _, c := range jc {
		cs = append(cs, c)
	}
	return json.Marshal(cs)
}

// DelegateRailNode is delegate of RailNode
type DelegateRailNode struct {
	DelegateNode
}

// B returns reference of Base Object
func (drn DelegateRailNode) B() *Base {
	return &drn.Base
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
	DelegateEdge

	From    *DelegateRailNode `json:"-"`
	To      *DelegateRailNode `json:"-"`
	Reverse *DelegateRailEdge `json:"-"`

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
