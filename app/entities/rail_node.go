package entities

import (
	"fmt"
	"time"
)

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Base
	Owner
	Point

	M            *Model             `gorm:"-" json:"-"`
	InEdges      map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdges     map[uint]*RailEdge `gorm:"-" json:"-"`
	OverPlatform *Platform          `gorm:"-" json:"-"`
	RailLines    map[uint]*RailLine `gorm:"-" json:"-"`
	InTasks      map[uint]*LineTask `gorm:"-" json:"-"`
	OutTasks     map[uint]*LineTask `gorm:"-" json:"-"`
	// key is id of RailNode
	Tracks map[uint]*Track `gorm:"-" json:"-"`

	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewRailNode create new instance.
func (m *Model) NewRailNode(o *Player, x float64, y float64) *RailNode {
	rn := &RailNode{
		Base:  NewBase(m.GenID(RAILNODE)),
		Point: NewPoint(x, y),
	}
	rn.Init(m)
	rn.Resolve(o)
	rn.Marshal()
	m.Add(rn)
	return rn
}

func (rn *RailNode) Extend(x float64, y float64) (*RailNode, *RailEdge) {
	to := rn.M.NewRailNode(rn.Own, x, y)
	e1 := rn.M.NewRailEdge(rn, to)
	e2 := rn.M.NewRailEdge(to, rn)

	e1.Resolve(e2)
	e2.Resolve(e1)

	rn.Own.ReRouting = true

	eachLineTask(rn.InTasks, func(lt *LineTask) {
		if lt.RailLine.AutoExt {
			lt.InsertRailEdge(e1)
		}
	})
	return to, e1
}

// Idx returns unique id field.
func (rn *RailNode) Idx() uint {
	return rn.ID
}

// Type returns type of entitiy
func (rn *RailNode) Type() ModelType {
	return RAILNODE
}

// Init makes map
func (rn *RailNode) Init(m *Model) {
	rn.M = m
	rn.InEdges = make(map[uint]*RailEdge)
	rn.OutEdges = make(map[uint]*RailEdge)
	rn.RailLines = make(map[uint]*RailLine)
	rn.InTasks = make(map[uint]*LineTask)
	rn.OutTasks = make(map[uint]*LineTask)
	rn.Tracks = make(map[uint]*Track)
}

// Pos returns location
func (rn *RailNode) Pos() *Point {
	return &rn.Point
}

// IsIn returns it should be view or not.
func (rn *RailNode) IsIn(x float64, y float64, scale float64) bool {
	return rn.Pos().IsIn(x, y, scale)
}

// Resolve set reference
func (rn *RailNode) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			rn.Owner = NewOwner(obj)
			obj.Resolve(rn)
		case *Platform:
			rn.OverPlatform = obj
		case *RailLine:
			rn.RailLines[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	rn.Marshal()
}

// Marshal set id from reference
func (rn *RailNode) Marshal() {
	if rn.OverPlatform != nil {
		rn.PlatformID = rn.OverPlatform.ID
	}
}

func (rn *RailNode) UnMarshal() {
	rn.Resolve(rn.M.Find(PLAYER, rn.OwnerID))
}

// BeforeDelete clear reference
func (rn *RailNode) BeforeDelete() {
	rn.Own.UnResolve(rn)
}

func (rn *RailNode) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Platform:
			rn.OverPlatform = nil
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Permits represents Player is permitted to control
func (rn *RailNode) Permits(o *Player) bool {
	return rn.Owner.Permits(o)
}

// CheckDelete checks remaining reference
func (rn *RailNode) CheckDelete() error {
	for _, re := range rn.OutEdges {
		if err := re.CheckDelete(); err != nil {
			return fmt.Errorf("blocked by OutEdges of %v (%v)", re, err)
		}
	}
	for _, re := range rn.InEdges {
		if err := re.CheckDelete(); err != nil {
			return fmt.Errorf("blocked by InEdges of %v (%v)", re, err)
		}
	}
	if rn.OverPlatform != nil {
		return fmt.Errorf("blocked by OverPlatform of %v", rn.OverPlatform)
	}
	return nil
}

func (rn *RailNode) Delete() {
	for _, re := range rn.OutEdges {
		re.Delete()
	}
	for _, re := range rn.InEdges {
		re.Delete()
	}
	rn.M.Delete(rn)
}

func (rn *RailNode) IsNew() bool {
	return rn.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (rn *RailNode) IsChanged(after ...time.Time) bool {
	return rn.Base.IsChanged(after...)
}

// Reset set status as not changed
func (rn *RailNode) Reset() {
	rn.Base.Reset()
}

// String represents status
func (rn *RailNode) String() string {
	rn.Marshal()
	ostr := ""
	if rn.Own != nil {
		ostr = fmt.Sprintf(":%s", rn.Own.Short())
	}
	pstr := ""
	if rn.OverPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", rn.OverPlatform.ID)
	}
	return fmt.Sprintf("%s(%d):i=%d,o=%d%s:%v%s", rn.Type().Short(),
		rn.ID, len(rn.InEdges), len(rn.OutEdges), pstr, rn.Pos(), ostr)
}
