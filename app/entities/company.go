package entities

// Company is the destination of Human
type Company struct {
	Model
	Point
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
		in:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
}

// Idx returns unique id field.
func (c *Company) Idx() uint {
	return c.ID
}

// Init creates map.
func (c *Company) Init() {
	c.Model.Init()
	c.Point.Init()
	c.in = make(map[uint]*Step)
	c.Targets = make(map[uint]*Human)
}

// Pos returns location.
func (c *Company) Pos() *Point {
	return &c.Point
}

// IsIn returns it should be view or not.
func (c *Company) IsIn(center *Point, scale float64) bool {
	return c.Pos().IsIn(center, scale)
}

// Out returns where it can go to.
func (c *Company) Out() map[uint]*Step {
	return nil
}

// In returns where it comes from.
func (c *Company) In() map[uint]*Step {
	return c.in
}

// ResolveRef do nothing (for implements Resolvable)
func (c *Company) ResolveRef() {
	// do-nothing
}
