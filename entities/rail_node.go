package entities

import (
	"fmt"
)

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Base
	Persistence
	Point

	InEdges  map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdges map[uint]*RailEdge `gorm:"-" json:"-"`
	// Tracks represents list of RailNode can be arrived at via specified OutEdge. key is id of RailEdge
	Tracks map[uint]map[uint]bool `gorm:"-" json:"-"`

	OverPlatform *Platform          `gorm:"-" json:"-"`
	InTasks      map[uint]*LineTask `gorm:"-" json:"-"`
	OutTasks     map[uint]*LineTask `gorm:"-" json:"-"`

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
	rn.Resolve(o)
	rn.Marshal()
	o.ReRouting = true
	m.Add(rn)
	return rn
}

// Extend returns new RailNode which position is (x, y)
func (rn *RailNode) Extend(x float64, y float64) (*RailNode, *RailEdge) {
	to := rn.M.NewRailNode(rn.O, x, y)
	return to, rn.Connect(to)
}

// Connect returns new RailEdge from rn to to
func (rn *RailNode) Connect(to *RailNode) *RailEdge {
	e1 := rn.M.NewRailEdge(rn, to)
	e2 := rn.M.NewRailEdge(to, rn)

	e1.Resolve(e2)
	e2.Resolve(e1)

	eachLineTask(rn.InTasks, func(lt *LineTask) {
		if lt.RailLine.AutoExt {
			lt.InsertRailEdge(e1)
		}
	})

	rn.O.ReRouting = true
	return e1
}

// B returns base information of this elements.
func (rn *RailNode) B() *Base {
	return &rn.Base
}

// P returns time information for database.
func (rn *RailNode) P() *Persistence {
	return &rn.Persistence
}

// Pos returns entities' position.
func (rn *RailNode) Pos() *Point {
	return &rn.Point
}

// Init makes map
func (rn *RailNode) Init(m *Model) {
	rn.Base.Init(RAILNODE, m)
	rn.InEdges = make(map[uint]*RailEdge)
	rn.OutEdges = make(map[uint]*RailEdge)
	rn.InTasks = make(map[uint]*LineTask)
	rn.OutTasks = make(map[uint]*LineTask)
	rn.Tracks = make(map[uint]map[uint]bool)
}

// Resolve set reference
func (rn *RailNode) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			rn.O = obj
			obj.Resolve(rn)
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
	if len(rn.InTasks)+len(rn.OutTasks) > 0 {
		return fmt.Errorf("blocked by LineTask (in=%d,out=%d)", len(rn.InTasks), len(rn.OutTasks))
	}
	return nil
}

// Delete removes this entity with related ones.
func (rn *RailNode) Delete() {
	if rn.OverPlatform != nil {
		rn.OverPlatform.Delete()
	}
	for _, re := range rn.OutEdges {
		re.Delete()
	}
	for _, re := range rn.InEdges {
		re.Delete()
	}
	rn.O.ReRouting = true
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
		rn.ID, len(rn.InEdges), len(rn.OutEdges), pstr, rn.Point, ostr)
}
