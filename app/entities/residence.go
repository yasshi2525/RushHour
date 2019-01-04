package entities

import (
	"fmt"
	"time"
)

// Residence generate Human in a period
type Residence struct {
	Base
	Owner
	Point
	out      map[uint]*Step
	Targets  map[uint]*Human `gorm:"-"        json:"-"`
	Capacity uint            `gorm:"not null" json:"capacity"`
	// Wait represents how msec after it generates Human
	Wait float64 `gorm:"not null" json:"wait"`
	Name string  `                json:"name"`
}

// NewResidence create new instance without setting parameters
func NewResidence(id uint, o *Player, x float64, y float64) *Residence {
	r := &Residence{
		Base:  NewBase(id),
		Owner: NewOwner(o),
		Point: NewPoint(x, y),
	}
	r.Init()
	return r
}

// Idx returns unique id field.
func (r *Residence) Idx() uint {
	return r.ID
}

// Type returns type of entitiy
func (r *Residence) Type() ModelType {
	return RESIDENCE
}

// Init creates map.
func (r *Residence) Init() {
	r.out = make(map[uint]*Step)
	r.Targets = make(map[uint]*Human)
}

// Pos returns location
func (r *Residence) Pos() *Point {
	return &r.Point
}

// IsIn returns it should be view or not.
func (r *Residence) IsIn(center *Point, scale float64) bool {
	return r.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (r *Residence) Out() map[uint]*Step {
	return r.out
}

// In returns where it comes from
func (r *Residence) In() map[uint]*Step {
	return nil
}

// Resolve set reference
func (r *Residence) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Human:
			r.Targets[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	r.ResolveRef()
}

// ResolveRef do nothing (for implements Resolvable)
func (r *Residence) ResolveRef() {
	// do-nothing
}

// UnRef remove reference of related entity
func (r *Residence) UnRef() {
}

// Permits represents Player is permitted to control
func (r *Residence) Permits(o *Player) bool {
	return o.Level == Admin
}

// CheckRemove check remaining reference
func (r *Residence) CheckRemove() error {
	return nil
}

// IsChanged returns true when it is changed after Backup()
func (r *Residence) IsChanged(after ...time.Time) bool {
	return r.Base.IsChanged(after)
}

// Reset set status as not changed
func (r *Residence) Reset() {
	r.Base.Reset()
}

func (r *Residence) String() string {
	return fmt.Sprintf("%s(%d):i=0,o=%d,h=%d:%v:%s", Meta.Attr[r.Type()].Short,
		r.ID, len(r.out), len(r.Targets), r.Pos(), r.Name)
}
