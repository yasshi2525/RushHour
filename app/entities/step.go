package entities

// Step represents two Relayable is logically connected.
// Step is out of target for persistence because it can derived by other resources.
type Step struct {
	ID     uint
	from   Relayable
	to     Relayable
	Weight float64
}

// NewStep create new instance and relation to Junction
func NewStep(id uint, f Relayable, t Relayable, weight float64) *Step {
	step := &Step{
		ID:     id,
		from:   f,
		to:     t,
		Weight: weight,
	}
	f.Out()[step.ID] = step
	t.In()[step.ID] = step
	return step
}

// Idx returns unique id field.
func (s *Step) Idx() uint {
	return s.ID
}

// Init do nothing
func (s *Step) Init() {
}

// Pos returns center
func (s *Step) Pos() *Point {
	return s.from.Pos().Center(s.to.Pos())
}

// IsIn returns it should be view or not.
func (s *Step) IsIn(center *Point, scale float64) bool {
	return s.from.Pos().IsInLine(s.to.Pos(), center, scale)
}

// From returns where Step comes from
func (s *Step) From() Relayable {
	return s.from
}

// To returns where Step goes to
func (s *Step) To() Relayable {
	return s.to
}

// Cost is calculated by distance * weight of Step.
func (s *Step) Cost() float64 {
	return s.from.Pos().Dist(s.to.Pos()) * s.Weight
}

// Unrelate delete selt from related Junction.
func (s *Step) Unrelate() {
	delete(s.from.Out(), s.ID)
	delete(s.to.In(), s.ID)
}
