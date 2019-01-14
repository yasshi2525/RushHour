package entities

import (
	"fmt"
	"math/rand"
	"time"
)

// Residence generate Human in a period
type Residence struct {
	Base
	Owner
	Point

	Capacity uint `gorm:"not null" json:"capacity"`
	// Wait represents how msec after it generates Human
	Wait float64 `gorm:"not null" json:"wait"`
	Name string  `                json:"name"`

	M       *Model          `gorm:"-" json:"-"`
	Targets map[uint]*Human `gorm:"-"        json:"-"`
	out     map[uint]*Step
}

// NewResidence create new instance without setting parameters
func (m *Model) NewResidence(o *Player, x float64, y float64) *Residence {
	r := &Residence{
		Base:     NewBase(m.GenID(RESIDENCE)),
		Owner:    NewOwner(o),
		Point:    NewPoint(x, y),
		Capacity: Const.Residence.Capacity,
		Wait:     Const.Residence.Interval.D.Seconds() * rand.Float64(),
	}
	r.Init(m)
	r.Resolve()
	r.Marshal()
	m.Add(r)

	r.GenOutSteps()
	return r
}

func (r *Residence) GenOutSteps() {
	// R -> C
	for _, c := range r.M.Companies {
		r.M.NewStep(r, c)
	}
	// R -> G
	for _, g := range r.M.Gates {
		r.M.NewStep(r, g)
	}
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
func (r *Residence) Init(m *Model) {
	r.M = m
	r.out = make(map[uint]*Step)
	r.Targets = make(map[uint]*Human)
}

// Pos returns location
func (r *Residence) Pos() *Point {
	return &r.Point
}

// IsIn returns it should be view or not.
func (r *Residence) IsIn(x float64, y float64, scale float64) bool {
	return r.Pos().IsIn(x, y, scale)
}

// OutSteps returns where it can go to
func (r *Residence) OutSteps() map[uint]*Step {
	return r.out
}

// InSteps returns where it comes from
func (r *Residence) InSteps() map[uint]*Step {
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
	r.Marshal()
}

// Marshal do nothing (for implements Resolvable)
func (r *Residence) Marshal() {
	// do-nothing
}

func (r *Residence) UnMarshal() {
	
}

// UnRef remove reference of related entity
func (r *Residence) UnRef() {
}

// Permits represents Player is permitted to control
func (r *Residence) Permits(o *Player) bool {
	return o.Level == Admin
}

// CheckDelete check remaining reference
func (r *Residence) CheckDelete() error {
	return nil
}

func (r *Residence) Delete() {
	for _, h := range r.Targets {
		r.M.Delete(h)
	}
	for _, s := range r.out {
		r.M.Delete(s)
	}
	r.M.Delete(r)
}

func (r *Residence) IsNew() bool {
	return r.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (r *Residence) IsChanged(after ...time.Time) bool {
	return r.Base.IsChanged(after...)
}

// Reset set status as not changed
func (r *Residence) Reset() {
	r.Base.Reset()
}

func (r *Residence) String() string {
	r.Marshal()
	return fmt.Sprintf("%s(%d):i=0,o=%d,h=%d:%v:%s", r.Type().Short(),
		r.ID, len(r.out), len(r.Targets), r.Pos(), r.Name)
}
