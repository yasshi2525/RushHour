package entities

// Step represents two Locationable is logically connected.
// Step is out of target for persistence because it can derived by other resources.
type Step struct {
	ID     uint
	from   Locationable
	to     Locationable
	Weight float64
}

// NewStep create new instance and relation to Junction
func NewStep(id uint, f Locationable, t Locationable, weight float64) *Step {
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

// From returns where Step comes from
func (s *Step) From() Locationable {
	return s.from
}

// To returns where Step goes to
func (s *Step) To() Locationable {
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
