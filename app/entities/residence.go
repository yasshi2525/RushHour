package entities

// Residence generate Human in a period
type Residence struct {
	Model
	Point
	out       map[uint]*Step
	in        map[uint]*Step
	Targets   map[uint]*Human `gorm:"-" json:"-"`
	Capacity  uint            `gorm:"not null" json:"capacity"`
	Available float64         `gorm:"not null" json:"available"`
}

// NewResidence create new instance without setting parameters
func NewResidence(id uint, x float64, y float64) *Residence {
	return &Residence{
		Model:   NewModel(id),
		Point:   NewPoint(x, y),
		out:     make(map[uint]*Step),
		in:      make(map[uint]*Step),
		Targets: make(map[uint]*Human),
	}
}

// Idx returns unique id field.
func (r *Residence) Idx() uint {
	return r.ID
}

// Pos returns location
func (r *Residence) Pos() *Point {
	return &r.Point
}

// Out returns where it can go to
func (r *Residence) Out() map[uint]*Step {
	return r.out
}

// In returns where it comes from
func (r *Residence) In() map[uint]*Step {
	return r.in
}

// IsIn returns it should be view or not.
func (r *Residence) IsIn(center *Point, scale float64) bool {
	return r.Pos().IsIn(center, scale)
}

// ResolveRef do nothing (for implements Resolvable)
func (r *Residence) ResolveRef() {
	// do-nothing
}
