package entities

import (
	"fmt"
)

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Base
	Persistence
	Shape
	Point

	InEdges      map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdges     map[uint]*RailEdge `gorm:"-" json:"-"`
	OverPlatform *Platform          `gorm:"-" json:"-"`
	InTasks      map[uint]*LineTask `gorm:"-" json:"-"`
	OutTasks     map[uint]*LineTask `gorm:"-" json:"-"`
	// key is id of RailNode
	Tracks map[uint]*Track `gorm:"-" json:"-"`

	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewRailNode create new instance.
func (m *Model) NewRailNode(o *Player, x float64, y float64) *RailNode {
	rn := &RailNode{
		Base:        m.NewBase(RAILNODE, o),
		Persistence: NewPersistence(),
		Point:       NewPoint(x, y),
	}
	rn.Init(m)
	m.Add(rn)
	return rn
}

// Extend returns new RailNode which position is (x, y)
func (rn *RailNode) Extend(x float64, y float64) (*RailNode, *RailEdge) {
	to := rn.M.NewRailNode(rn.O, x, y)
	e1 := rn.M.NewRailEdge(rn, to)
	e2 := rn.M.NewRailEdge(to, rn)

	e1.Resolve(e2)
	e2.Resolve(e1)

	eachLineTask(rn.InTasks, func(lt *LineTask) {
		if lt.RailLine.AutoExt {
			lt.InsertRailEdge(e1)
		}
	})
	return to, e1
}

// B returns base information of this elements.
func (rn *RailNode) B() *Base {
	return &rn.Base
}

// P returns time information for database.
func (rn *RailNode) P() *Persistence {
	return &rn.Persistence
}

// S returns entities' position.
func (rn *RailNode) S() *Shape {
	return &rn.Shape
}

// Init makes map
func (rn *RailNode) Init(m *Model) {
	rn.Base.Init(RAILNODE, m)
	rn.Shape.P1 = &rn.Point
	rn.InEdges = make(map[uint]*RailEdge)
	rn.OutEdges = make(map[uint]*RailEdge)
	rn.InTasks = make(map[uint]*LineTask)
	rn.OutTasks = make(map[uint]*LineTask)
	rn.Tracks = make(map[uint]*Track)
}

// Resolve set reference
func (rn *RailNode) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			rn.O = obj
		case *Platform:
			rn.OverPlatform = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	rn.Marshal()
}

// Marshal set id from reference
func (rn *RailNode) Marshal() {
	if rn.O != nil {
		rn.OwnerID = rn.O.ID
	}
	if rn.OverPlatform != nil {
		rn.PlatformID = rn.OverPlatform.ID
	}
}

// UnMarshal set reference from id.
func (rn *RailNode) UnMarshal() {
	rn.Resolve(rn.M.Find(PLAYER, rn.OwnerID))
}

// BeforeDelete clear reference
func (rn *RailNode) BeforeDelete() {
	rn.O.UnResolve(rn)
}

// UnResolve unregisters specified refernce.
func (rn *RailNode) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Platform:
			rn.OverPlatform = nil
			rn.PlatformID = ZERO
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// CheckDelete checks remaining reference
func (rn *RailNode) CheckDelete() error {
	if rn.OverPlatform != nil {
		return fmt.Errorf("blocked by OverPlatform of %v", rn.OverPlatform)
	}
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
	return nil
}

// Delete removes this entity with related ones.
func (rn *RailNode) Delete(force bool) {
	for _, re := range rn.OutEdges {
		re.Delete(false)
	}
	for _, re := range rn.InEdges {
		re.Delete(false)
	}
	rn.M.Delete(rn)
}

// String represents status
func (rn *RailNode) String() string {
	rn.Marshal()
	ostr := ""
	if rn.O != nil {
		ostr = fmt.Sprintf(":%s", rn.O.Short())
	}
	pstr := ""
	if rn.OverPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", rn.OverPlatform.ID)
	}
	return fmt.Sprintf("%s(%d):i=%d,o=%d%s:%v%s", rn.Type().Short(),
		rn.ID, len(rn.InEdges), len(rn.OutEdges), pstr, rn.Pos(), ostr)
}
