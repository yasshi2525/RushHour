package entities

import (
	"fmt"
	"time"
)

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Base
	Owner

	outStep map[uint]*Step
	inStep  map[uint]*Step

	InStation  *Station           `gorm:"-" json:"-"`
	WithGate   *Gate              `gorm:"-" json:"-"`
	OnRailNode *RailNode          `gorm:"-" json:"-"`
	Passengers map[uint]*Human    `gorm:"-" json:"-"`
	LineTasks  map[uint]*LineTask `gorm:"-" json:"-"`

	Trains map[uint]*Train `gorm:"-" json:"-"`

	Capacity uint `gorm:"not null" json:"cap"`
	Occupied uint `gorm:"-"        json:"used"`

	StationID  uint `gorm:"not null" json:"stid"`
	GateID     uint `gorm:"-"        json:"gid"`
	RailNodeID uint `gorm:"not null" json:"rnid"`
}

// NewPlatform creates instance
func NewPlatform(pid uint, rn *RailNode, g *Gate, st *Station) *Platform {
	p := &Platform{
		Base:       NewBase(pid),
		Owner:      rn.Owner,
		OnRailNode: rn,
		InStation:  st,
		WithGate:   g,
	}
	p.Init()
	p.ResolveRef()
	rn.Resolve(p)
	g.Resolve(p)
	st.Resolve(p)
	return p
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
	p.outStep = make(map[uint]*Step)
	p.inStep = make(map[uint]*Step)
	p.Passengers = make(map[uint]*Human)
	p.Trains = make(map[uint]*Train)
	p.LineTasks = make(map[uint]*LineTask)
}

// Pos returns location
func (p *Platform) Pos() *Point {
	if p.OnRailNode == nil {
		return nil
	}
	return p.OnRailNode.Pos()
}

// IsIn returns it should be view or not.
func (p *Platform) IsIn(x float64, y float64, scale float64) bool {
	return p.Pos().IsIn(x, y, scale)
}

// OutStep returns where it can go to
func (p *Platform) OutStep() map[uint]*Step {
	return p.outStep
}

// InStep returns where it comes from
func (p *Platform) InStep() map[uint]*Step {
	return p.inStep
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
			obj.Resolve(p)
		case *LineTask:
			p.LineTasks[obj.ID] = obj
		case *Train:
			p.Trains[obj.ID] = obj
			obj.Resolve(p)
		case *Human:
			p.Passengers[obj.ID] = obj
			p.Occupied++
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	p.ResolveRef()
}

// ResolveRef set id from reference
func (p *Platform) ResolveRef() {
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

// CheckRemove checks related reference
func (p *Platform) CheckRemove() error {
	if len(p.LineTasks) > 0 {
		return fmt.Errorf("blocked by LineTask of %v", p.Trains)
	}
	if len(p.Trains) > 0 {
		return fmt.Errorf("blocked by Train of %v", p.Trains)
	}
	return nil
}

// UnRef delete related reference.
func (p *Platform) UnRef() {
	for _, h := range p.Passengers {
		h.SetOnPlatform(nil)
		delete(p.Passengers, h.ID)
	}
	for _, t := range p.Trains {
		t.OnPlatform = nil
		delete(p.Trains, t.ID)
	}
	for _, lt := range p.LineTasks {
		lt.UnRef()
		delete(p.LineTasks, lt.ID)
	}
}

// Permits represents Player is permitted to control
func (p *Platform) Permits(o *Player) bool {
	return p.Owner.Permits(o)
}

// IsChanged returns true when it is changed after Backup()
func (p *Platform) IsChanged(after ...time.Time) bool {
	return p.Base.IsChanged(after...)
}

// Reset set status as not changed
func (p *Platform) Reset() {
	p.Base.Reset()
}

// String represents status
func (p *Platform) String() string {
	p.ResolveRef()
	ostr := ""
	if p.Own != nil {
		ostr = fmt.Sprintf(":%s", p.Own.Short())
	}
	posstr := ""
	if p.Pos() != nil {
		posstr = fmt.Sprintf(":%s", p.Pos())
	}
	nmstr := ""
	if p.InStation != nil {
		nmstr = fmt.Sprintf(":%s", p.InStation.Name)
	}
	return fmt.Sprintf("%s(%d):st=%d,g=%d,rn=%d,i=%d,o=%d,h=%d/%d%s%s%s",
		Meta.Attr[p.Type()].Short,
		p.ID, p.StationID, p.GateID, p.RailNodeID,
		len(p.inStep), len(p.outStep), len(p.Passengers), p.Capacity,
		posstr, ostr, nmstr)
}
