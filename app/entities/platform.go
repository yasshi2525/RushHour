package entities

import (
	"fmt"
)

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Base
	Owner

	out map[uint]*Step
	in  map[uint]*Step

	InStation  *Station        `gorm:"-" json:"-"`
	WithGate   *Gate           `gorm:"-" json:"-"`
	OnRailNode *RailNode       `gorm:"-" json:"-"`
	Passenger  map[uint]*Human `gorm:"-" json:"-"`

	Trains map[uint]*Train `gorm:"-" json:"-"`

	Capacity uint `gorm:"not null" json:"cap"`
	Occupied uint `gorm:"-"        json:"used"`

	StationID  uint `gorm:"not null" json:"stid"`
	GateID     uint `gorm:"-" json:"gid"`
	RailNodeID uint `gorm:"not null" json:"rnid"`
}

// Idx returns unique id field.
func (p *Platform) Idx() uint {
	return p.ID
}

// Type returns type of entitiy
func (p *Platform) Type() ModelType {
	return PLATFORM
}

// Init creates map.
func (p *Platform) Init() {
	p.out = make(map[uint]*Step)
	p.in = make(map[uint]*Step)
	p.Passenger = make(map[uint]*Human)
	p.Trains = make(map[uint]*Train)
}

// Pos returns location
func (p *Platform) Pos() *Point {
	return p.OnRailNode.Pos()
}

// IsIn returns it should be view or not.
func (p *Platform) IsIn(center *Point, scale float64) bool {
	return p.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (p *Platform) Out() map[uint]*Step {
	return p.out
}

// In returns where it comes from
func (p *Platform) In() map[uint]*Step {
	return p.in
}

// Resolve set reference
func (p *Platform) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			p.Owner, p.OnRailNode = obj.Owner, obj
			obj.Resolve(p)
		case *Station:
			p.InStation = obj
			obj.Resolve(p)
		case *Gate:
			p.WithGate = obj
		case *Train:
			p.Trains[obj.ID] = obj
			obj.Resolve(p)
		case *Human:
			p.Passenger[obj.ID] = obj
			p.Occupied++
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	p.ResolveRef()
}

// ResolveRef set id from reference
func (p *Platform) ResolveRef() {
	p.Owner.ResolveRef()
	if p.OnRailNode != nil {
		p.RailNodeID = p.OnRailNode.ID
	}
	if p.WithGate != nil {
		p.GateID = p.WithGate.ID
	}
	if p.InStation != nil {
		p.StationID = p.InStation.ID
	}
}

// Permits represents Player is permitted to control
func (p *Platform) Permits(o *Player) bool {
	return p.Owner.Permits(o)
}

// String represents status
func (p *Platform) String() string {
	return fmt.Sprintf("%s(%d):st=%d,g=%d,i=%d,o=%d,h=%d/%d:%v:%s",
		Meta.Attr[p.Type()].Short,
		p.ID, p.InStation.ID, p.WithGate.ID,
		len(p.in), len(p.out), len(p.Passenger), p.Capacity,
		p.Pos(), p.InStation.Name)
}
