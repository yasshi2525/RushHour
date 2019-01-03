package entities

import (
	"fmt"
)

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Base
	Owner
	Point
	InEdge  map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdge map[uint]*RailEdge `gorm:"-" json:"-"`

	OverPlatform *Platform `gorm:"-" json:"-"`
	PlatformID   uint      `gorm:"-" json:"pid,omitempty"`
}

// NewRailNode create new instance.
func NewRailNode(id uint, o *Player, x float64, y float64) *RailNode {
	re := &RailNode{
		Base:    NewBase(id),
		Owner:   NewOwner(o),
		Point:   NewPoint(x, y),
	}
	re.Init()
	return re
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
func (rn *RailNode) Init() {
	rn.InEdge = make(map[uint]*RailEdge)
	rn.OutEdge = make(map[uint]*RailEdge)
}

// Pos returns location
func (rn *RailNode) Pos() *Point {
	return &rn.Point
}

// IsIn returns it should be view or not.
func (rn *RailNode) IsIn(center *Point, scale float64) bool {
	return rn.Pos().IsIn(center, scale)
}

// Resolve set reference
func (rn *RailNode) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			rn.Owner = NewOwner(obj)
		case *Platform:
			rn.OverPlatform = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	rn.ResolveRef()
}

// ResolveRef set id from reference
func (rn *RailNode) ResolveRef() {
	if rn.OverPlatform != nil {
		rn.PlatformID = rn.OverPlatform.ID
	}
}

// UnRef clear reference
func (rn *RailNode) UnRef() {
	// do nothing
}

// Permits represents Player is permitted to control
func (rn *RailNode) Permits(o *Player) bool {
	return rn.Owner.Permits(o)
}

// CheckRemove checks remaining reference
func (rn *RailNode) CheckRemove() error {
	if len(rn.InEdge) > 0 {
		return fmt.Errorf("blocked by InEdge of %v", rn.InEdge)
	}
	if len(rn.OutEdge) > 0 {
		return fmt.Errorf("blocked by OutEdge of %v", rn.OutEdge)
	}
	if rn.OverPlatform != nil {
		return fmt.Errorf("blocked by Platform of %v", rn.OverPlatform)
	}
	return nil
}

// IsChanged returns true when it is changed after Backup()
func (rn *RailNode) IsChanged() bool {
	return rn.Base.IsChanged()
}

// Reset set status as not changed
func (rn *RailNode) Reset() {
	rn.Base.Reset()
}

// String represents status
func (rn *RailNode) String() string {
	ostr := ""
	if rn.Own != nil {
		ostr = fmt.Sprintf(":%s", rn.Own.Short())
	}
	pstr := ""
	if rn.OverPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", rn.OverPlatform.ID)
	}
	return fmt.Sprintf("%s(%d):i=%d,o=%d%s:%v%s", Meta.Attr[rn.Type()].Short,
		rn.ID, len(rn.InEdge), len(rn.OutEdge), pstr, rn.Pos(), ostr)
}
