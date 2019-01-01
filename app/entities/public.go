package entities

// Company is the destination of Human
type Company struct {
	Model
	Loc     *Point
	Out     map[uint]*Step  `gorm:"-" json:"-"`
	In      map[uint]*Step  `gorm:"-" json:"-"`
	Targets map[uint]*Human `gorm:"-" json:"-"`
	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64 `gorm:"not null" json:"scale"`
}

// NewCompany create new instance without setting parameters
func NewCompany(id uint, x float64, y float64) *Company {
	return &Company{
		Model:   NewModel(id),
		Loc:     NewPoint(x, y),
		Out:     make(map[uint]*Step),
		In:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
}

// Pos returns location
func (c *Company) Pos() *Point {
	return c.Loc
}

// OutStep returns where it can go to
func (c *Company) OutStep() map[uint]*Step {
	return c.Out
}

// InStep returns where it comes from
func (c *Company) InStep() map[uint]*Step {
	return c.In
}

// ResolveRef do nothing (for implements Resolvable)
func (c *Company) ResolveRef() {
	// do-nothing
}

// Residence generate Human in a period
type Residence struct {
	Model
	Loc       *Point
	Out       map[uint]*Step  `gorm:"-" json:"-"`
	In        map[uint]*Step  `gorm:"-" json:"-"`
	Targets   map[uint]*Human `gorm:"-" json:"-"`
	Capacity  uint            `gorm:"not null" json:"capacity"`
	Available float64         `gorm:"not null" json:"available"`
}

// NewResidence create new instance without setting parameters
func NewResidence(id uint, x float64, y float64) *Residence {
	return &Residence{
		Model:   NewModel(id),
		Loc:     NewPoint(x, y),
		Out:     make(map[uint]*Step),
		In:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
}

// Pos returns location
func (r *Residence) Pos() *Point {
	return r.Loc
}

// OutStep returns where it can go to
func (r *Residence) OutStep() map[uint]*Step {
	return r.Out
}

// InStep returns where it comes from
func (r *Residence) InStep() map[uint]*Step {
	return r.In
}

// ResolveRef do nothing (for implements Resolvable)
func (r *Residence) ResolveRef() {
	// do-nothing
}
