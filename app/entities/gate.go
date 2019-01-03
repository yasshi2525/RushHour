package entities

import (
	"fmt"
)

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	Base
	Owner

	out map[uint]*Step
	in  map[uint]*Step

	InStation    *Station  `gorm:"-" json:"-"`
	WithPlatform *Platform `gorm:"-" json:"-"`

	// Num represents how many Human can pass at the same time
	Num uint `gorm:"not null" json:"num"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null" json:"mobility"`
	// Occupied represents how many Gate are used by Human.
	Occupied uint `gorm:"not null" json:"occupied"`

	StationID  uint `gorm:"not null" json:"stid"`
	PlatformID uint `gorm:"-"        json:"pid"`
}

// Idx returns unique id field.
func (g *Gate) Idx() uint {
	return g.ID
}

// Type returns type of entitiy
func (g *Gate) Type() ModelType {
	return GATE
}

// Init creates map.
func (g *Gate) Init() {
	g.out = make(map[uint]*Step)
	g.in = make(map[uint]*Step)
}

// Pos returns location
func (g *Gate) Pos() *Point {
	return g.WithPlatform.Pos()
}

// IsIn returns it should be view or not.
func (g *Gate) IsIn(center *Point, scale float64) bool {
	return g.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (g *Gate) Out() map[uint]*Step {
	return g.out
}

// In returns where it comes from
func (g *Gate) In() map[uint]*Step {
	return g.in
}

// Resolve set reference
func (g *Gate) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			g.Own = obj
		case *Station:
			g.InStation = obj
			obj.Resolve(g)
		case *Platform:
			g.WithPlatform = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	g.ResolveRef()
}

// ResolveRef set id from reference
func (g *Gate) ResolveRef() {
	g.Owner.ResolveRef()
	if g.InStation != nil {
		g.StationID = g.InStation.ID
	}
	if g.WithPlatform != nil {
		g.PlatformID = g.WithPlatform.ID
	}
}

// Permits represents Player is permitted to control
func (g *Gate) Permits(o *Player) bool {
	return g.Owner.Permits(o)
}

// String represents status
func (g *Gate) String() string {
	return fmt.Sprintf("%s(%d):st=%d,p=%d,i=%d,o=%d:%v:%s", Meta.Attr[g.Type()].Short,
		g.ID, g.InStation.ID, g.WithPlatform.ID,
		len(g.in), len(g.out),
		g.Pos(), g.InStation.Name)
}
