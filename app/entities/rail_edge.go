package entities

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Model
	Owner
	from *RailNode
	to   *RailNode

	FromID uint `gorm:"not null"`
	ToID   uint `gorm:"not null"`
}

// NewRailEdge create new instance and relates RailNode
func NewRailEdge(id uint, f *RailNode, t *RailNode) *RailEdge {
	re := &RailEdge{
		Model: NewModel(id),
		Owner: f.Owner,
		from:  f,
		to:    t,
	}
	re.ResolveRef()

	f.OutEdge[re.ID] = re
	t.InEdge[re.ID] = re
	return re
}

// Idx returns unique id field.
func (re *RailEdge) Idx() uint {
	return re.ID
}

// Pos returns location
func (re *RailEdge) Pos() *Point {
	return re.from.Pos().Center(re.to.Pos())
}

// Out returns where it can go to
func (re *RailEdge) Out() map[uint]*Step {
	return nil
}

// In returns where it comes from
func (re *RailEdge) In() map[uint]*Step {
	return nil
}

// IsIn return true when from, to, center is in,
func (re *RailEdge) IsIn(center *Point, scale float64) bool {
	return re.from.Pos().IsIn(center, scale) ||
		re.to.Pos().IsIn(center, scale) ||
		re.Pos().IsIn(center, scale)
}

// From represents start point
func (re *RailEdge) From() Locationable {
	return re.from
}

// To represents end point
func (re *RailEdge) To() Locationable {
	return re.to
}

// Cost represents distance
func (re *RailEdge) Cost() float64 {
	return re.from.Pos().Dist(re.to.Pos())
}

// Unrelate delete relations to RailNode
func (re *RailEdge) Unrelate() {
	delete(re.from.OutEdge, re.ID)
	delete(re.to.InEdge, re.ID)
}

// Resolve set reference.
func (re *RailEdge) Resolve(from *RailNode, to *RailNode) {
	re.Owner, re.from, re.to = from.Owner, from, to

	from.OutEdge[re.ID] = re
	to.InEdge[re.ID] = re

	re.ResolveRef()
}

// ResolveRef set id from reference
func (re *RailEdge) ResolveRef() {
	re.Owner.ResolveRef()
	re.FromID = re.from.ID
	re.ToID = re.to.ID
}

// Permits represents Player is permitted to control
func (re *RailEdge) Permits(o *Player) bool {
	return re.Owner.Permits(o)
}
