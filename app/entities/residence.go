package entities

// Residence generate Human in a period
type Residence struct {
	Model
	Point
	out      map[uint]*Step
	in       map[uint]*Step
	Targets  map[uint]*Human `gorm:"-" json:"-"`
	Capacity uint            `gorm:"not null" json:"cap"`
	// Wait represents how msec after it generates Human
	Wait float64 `gorm:"not null" json:"wait"`
}

// NewResidence create new instance without setting parameters
func NewResidence(id uint, x float64, y float64) *Residence {
	r := &Residence{
		Model:   NewModel(id),
		Point:   NewPoint(x, y),
		out:     make(map[uint]*Step),
		in:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
	return r
}

// Idx returns unique id field.
func (r *Residence) Idx() uint {
	return r.ID
}

// Init creates map.
func (r *Residence) Init() {
	r.Model.Init()
	r.out = make(map[uint]*Step)
	r.in = make(map[uint]*Step)
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
	return r.in
}

// ResolveRef do nothing (for implements Resolvable)
func (r *Residence) ResolveRef() {
	// do-nothing
}
