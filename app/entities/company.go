package entities

import (
	"fmt"
)

// Company is the destination of Human
type Company struct {
	Base
	Persistence
	Point
	Shape

	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64 `gorm:"not null" json:"scale"`
	Name  string  `json:"name"`

	Targets map[uint]*Human `gorm:"-" json:"-"`
	in      map[uint]*Step
}

// NewCompany create new instance without setting parameters
func (m *Model) NewCompany(x float64, y float64) *Company {
	c := &Company{
		Base:        m.NewBase(COMPANY),
		Persistence: NewPersistence(),
		Point:       NewPoint(x, y),
		Scale:       Const.Company.Scale,
	}
	c.Init(m)
	c.Resolve()
	c.Marshal()
	m.Add(c)

	c.GenInSteps()
	return c
}

// B returns base information of this elements.
func (c *Company) B() *Base {
	return &c.Base
}

// P returns time information for database.
func (c *Company) P() *Persistence {
	return &c.Persistence
}

// S returns entities' position.
func (c *Company) S() *Shape {
	return &c.Shape
}

// GenInSteps generates and registers Step for this Company.
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

// Init creates map.
func (c *Company) Init(m *Model) {
	c.Base.Init(COMPANY, m)
	c.Shape.P1 = &c.Point
	c.in = make(map[uint]*Step)
	c.Targets = make(map[uint]*Human)
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
func (c *Company) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Human:
			c.Targets[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	c.Marshal()
}

// Marshal do nothing (for implements Resolvable)
func (c *Company) Marshal() {
	// do-nothing
}

// UnMarshal do nothing.
func (c *Company) UnMarshal() {

}

// BeforeDelete remove reference of related entity
func (c *Company) BeforeDelete() {
}

// CheckDelete check remaining reference.
func (c *Company) CheckDelete() error {
	return nil
}

// Delete removes this entity with related ones.
func (c *Company) Delete(force bool) {
	for _, h := range c.Targets {
		c.M.Delete(h)
	}
	for _, s := range c.in {
		c.M.Delete(s)
	}
	c.M.Delete(c)
}

// String represents status
func (c *Company) String() string {
	c.Marshal()
	return fmt.Sprintf("%s(%d):i=%d,o=0,h=%d:%v:%s", c.Type().Short(),
		c.ID, len(c.in), len(c.Targets), c.Pos(), c.Name)
}
