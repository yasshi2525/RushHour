package entities

// Company is the destination of Human
type Company struct {
	Model
	Point
	out     map[uint]*Step
	in      map[uint]*Step
	Targets map[uint]*Human `gorm:"-" json:"-"`
	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64 `gorm:"not null" json:"scale"`
}

// NewCompany create new instance without setting parameters
func NewCompany(id uint, x float64, y float64) *Company {
	return &Company{
		Model:   NewModel(id),
		Point:   NewPoint(x, y),
		out:     make(map[uint]*Step),
		in:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
}

// Idx returns unique id field.
func (c *Company) Idx() uint {
	return c.ID
}

// Pos returns location.
func (c *Company) Pos() *Point {
	return &c.Point
}

// Out returns where it can go to.
func (c *Company) Out() map[uint]*Step {
	return c.out
}

// In returns where it comes from.
func (c *Company) In() map[uint]*Step {
	return c.in
}

// IsIn returns it should be view or not.
func (c *Company) IsIn(center *Point, scale float64) bool {
	return c.Pos().IsIn(center, scale)
}

// ResolveRef do nothing (for implements Resolvable)
func (c *Company) ResolveRef() {
	// do-nothing
}
