package entities

import (
	"fmt"
	"time"
)

// Company is the destination of Human
type Company struct {
	Base
	Owner
	Point

	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64 `gorm:"not null" json:"scale"`
	Name  string  `json:"name"`

	M       *Model          `gorm:"-" json:"-"`
	Targets map[uint]*Human `gorm:"-" json:"-"`
	in      map[uint]*Step
}

// NewCompany create new instance without setting parameters
func (m *Model) NewCompany(o *Player, x float64, y float64) *Company {
	c := &Company{
		Base:  NewBase(m.GenID(COMPANY)),
		Owner: NewOwner(o),
		Point: NewPoint(x, y),
		Scale: Const.Company.Scale,
	}
	c.Init(m)
	c.Resolve()
	c.ResolveRef()
	m.Add(c)

	c.GenInSteps()
	return c
}

func (c *Company) GenInSteps() {
	// R -> C
	for _, r := range c.M.Residences {
		c.M.NewStep(r, c)
	}
	// G -> C
	for _, g := range c.M.Gates {
		c.M.NewStep(g, c)
	}
}

// Idx returns unique id field.
func (c *Company) Idx() uint {
	return c.ID
}

// Type returns type of entitiy
func (c *Company) Type() ModelType {
	return COMPANY
}

// Init creates map.
func (c *Company) Init(m *Model) {
	c.M = m
	c.in = make(map[uint]*Step)
	c.Targets = make(map[uint]*Human)
}

// Pos returns location.
func (c *Company) Pos() *Point {
	return &c.Point
}

// IsIn returns it should be view or not.
func (c *Company) IsIn(x float64, y float64, scale float64) bool {
	return c.Pos().IsIn(x, y, scale)
}

// OutSteps returns where it can go to.
func (c *Company) OutSteps() map[uint]*Step {
	return nil
}

// InSteps returns where it comes from.
func (c *Company) InSteps() map[uint]*Step {
	return c.in
}

// Resolve set reference
func (c *Company) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Human:
			c.Targets[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	c.ResolveRef()
}

// ResolveRef do nothing (for implements Resolvable)
func (c *Company) ResolveRef() {
	// do-nothing
}

// UnRef remove reference of related entity
func (c *Company) UnRef() {
}

// Permits represents Player is permitted to control
func (c *Company) Permits(o *Player) bool {
	return o.Level == Admin
}

// CheckDelete check remaining reference
func (c *Company) CheckDelete() error {
	return nil
}

func (c *Company) Delete() {
	for _, h := range c.Targets {
		c.M.Delete(h)
	}
	for _, s := range c.in {
		c.M.Delete(s)
	}
	c.M.Delete(c)
}

func (c *Company) IsNew() bool {
	return c.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (c *Company) IsChanged(after ...time.Time) bool {
	return c.Base.IsChanged(after...)
}

// Reset set status as not changed
func (c *Company) Reset() {
	c.Base.Reset()
}

// String represents status
func (c *Company) String() string {
	c.ResolveRef()
	return fmt.Sprintf("%s(%d):i=%d,o=0,h=%d:%v:%s", c.Type().Short(),
		c.ID, len(c.in), len(c.Targets), c.Pos(), c.Name)
}
