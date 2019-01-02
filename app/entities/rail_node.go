package entities

import (
	"fmt"
)

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Model
	Owner
	Point
	InEdge  map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdge map[uint]*RailEdge `gorm:"-" json:"-"`

	OverPlatform *Platform `gorm:"-" json:"-"`
	PlatformID   uint      `gorm:"-" json:"pid,omitempty"`
}

// NewRailNode create new instance.
func NewRailNode(id uint, o *Player, x float64, y float64) *RailNode {
	return &RailNode{
		Model:   NewModel(id),
		Owner:   NewOwner(o),
		Point:   NewPoint(x, y),
		InEdge:  make(map[uint]*RailEdge),
		OutEdge: make(map[uint]*RailEdge),
	}
}

// Idx returns unique id field.
func (rn *RailNode) Idx() uint {
	return rn.ID
}

// Init makes map
func (rn *RailNode) Init() {
	rn.Model.Init()
	rn.Owner.Init()
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
			rn.Own = obj
		case *Platform:
			rn.OverPlatform = obj
			obj.Resolve(rn)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	rn.ResolveRef()
}

// ResolveRef set id from reference
func (rn *RailNode) ResolveRef() {
	rn.Owner.ResolveRef()
	if rn.OverPlatform != nil {
		rn.PlatformID = rn.OverPlatform.ID
	}
}

// Permits represents Player is permitted to control
func (rn *RailNode) Permits(o *Player) bool {
	return rn.Owner.Permits(o)
}

// String represents status
func (rn *RailNode) String() string {
	pstr := ""
	if rn.OverPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", rn.OverPlatform.ID)
	}
	return fmt.Sprintf("%s(%d):i=%d,o=%d%s:%v", Meta.Static[RAILNODE].Short,
		rn.ID, len(rn.InEdge), len(rn.OutEdge), pstr, rn.Pos())
}
