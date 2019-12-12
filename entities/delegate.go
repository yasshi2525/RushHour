package entities

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DelegateMap represents Map for client view.
type DelegateMap struct {
	// Residences is the list of delegated Residence information
	Residences map[uint]*DelegateResidence `json:"residences"`
	// Companies is the list of delegated Company information
	Companies map[uint]*DelegateCompany `json:"companies"`
	// RailNodes is the list of delegated RailNode information
	RailNodes map[uint]*DelegateRailNode `json:"rail_nodes"`
	// RailEdges is the list of delegated RailEdge information
	RailEdges map[uint]*DelegateRailEdge `json:"rail_edges"`

	Values map[ModelType]reflect.Value `json:"-"`

	Timestamp int64 `json:"timestamp"`
}

// Init initialize map
func (dm *DelegateMap) Init() {
	dm.Residences = make(map[uint]*DelegateResidence)
	dm.Companies = make(map[uint]*DelegateCompany)
	dm.RailNodes = make(map[uint]*DelegateRailNode)
	dm.RailEdges = make(map[uint]*DelegateRailEdge)

	dm.Values = make(map[ModelType]reflect.Value)
	v := reflect.ValueOf(dm).Elem()
	for idx, ty := range []ModelType{RESIDENCE, COMPANY, RAILNODE, RAILEDGE} {
		dm.Values[ty] = v.Field(idx)
	}

	dm.Timestamp = time.Now().Unix()
}

// Add add delegatable object to map
func (dm *DelegateMap) Add(obj delegateLocalable) {
	if reflect.ValueOf(obj).IsNil() {
		return
	}
	dm.Values[obj.B().Type()].SetMapIndex(
		reflect.ValueOf(obj.B().Idx()), reflect.ValueOf(obj))
}

type delegateLocalable interface {
	B() *Base
}

type delegatePoint struct {
	SumX int     `json:"-"`
	SumY int     `json:"-"`
	AveX float64 `json:"x"`
	AveY float64 `json:"y"`
}

// DelegateNode represents delegation Localable.
type DelegateNode struct {
	Base

	Scale     int                `json:"-"`
	Pos       *delegatePoint     `json:"pos"`
	Multi     int                `json:"mul"`
	ParentIDs []uint             `json:"pids,omitempty"`
	ChildID   uint               `json:"cid,omitempty"`
	List      map[uint]Localable `json:"-"`
}

// NewDelegateNode creates instance.
func (ch *Chunk) NewDelegateNode(obj Localable, pids []uint) DelegateNode {
	scale := ch.Scale - ch.M.conf.MinScale
	x, y := Logarithm(obj.Pos().X, scale), Logarithm(obj.Pos().Y, scale)
	return DelegateNode{
		Base:      ch.M.NewBase(obj.B().T, obj.B().O),
		Scale:     ch.Scale,
		Pos:       &delegatePoint{x, y, float64(x), float64(y)},
		ParentIDs: pids,
		List:      make(map[uint]Localable),
	}
}

// Add accepts new instance ant increment count variable
func (dn *DelegateNode) Add(obj Localable) {
	dn.List[obj.B().ID] = obj
	dn.updateMulti()
	scale := dn.Scale - dn.M.conf.MinScale
	dn.Pos.SumX += Logarithm(obj.Pos().X, scale)
	dn.Pos.SumY += Logarithm(obj.Pos().Y, scale)
	dn.updatePos()
}

// Remove delete argument from list ant decrement count variable
func (dn *DelegateNode) Remove(obj Localable) {
	scale := dn.Scale - dn.M.conf.MinScale
	dn.Pos.SumX -= Logarithm(obj.Pos().X, scale)
	dn.Pos.SumY -= Logarithm(obj.Pos().Y, scale)
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
	if len(dn.List) > 0 {
		dn.Pos.AveX = float64(dn.Pos.SumX) / float64(len(dn.List))
		dn.Pos.AveY = float64(dn.Pos.SumY) / float64(len(dn.List))
	} else {
		dn.Pos.AveX = 0
		dn.Pos.AveY = 0
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

// DelegateCompany is delegate of RailNode
type DelegateCompany struct {
	DelegateNode
}

// B returns reference of Base Object
func (dc DelegateCompany) B() *Base {
	return &dc.Base
}

// DelegateRailNode is delegate of RailNode
type DelegateRailNode struct {
	DelegateNode
}

// B returns reference of Base Object
func (drn DelegateRailNode) B() *Base {
	return &drn.Base
}

// DelegateRailEdge is delegate of RailEdge
type DelegateRailEdge struct {
	DelegateEdge

	Reverse   *DelegateRailEdge `json:"-"`
	ReverseID uint              `json:"eid"`
}

// B returns reference of Base Object
func (dre DelegateRailEdge) B() *Base {
	return &dre.Base
}
