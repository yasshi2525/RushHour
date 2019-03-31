package entities

import (
	"fmt"
)

// Step represents two Relayable is logically connected.
// Step is out of target for persistence because it can derived by other resources.
type Step struct {
	Base
	Shape
	FromNode Relayable
	ToNode   Relayable
}

// NewStep create new instance and relation to Relayable
func (m *Model) NewStep(f Relayable, t Relayable) *Step {
	s := &Step{
		Base:     m.NewBase(STEP),
		Shape:    NewShapeEdge(f.S().Pos(), t.S().Pos()),
		FromNode: f,
		ToNode:   t,
	}
	s.Init(m)
	f.OutSteps()[s.ID] = s
	t.InSteps()[s.ID] = s
	m.Add(s)
	return s
}

// B returns base information of this elements.
func (s *Step) B() *Base {
	return &s.Base
}

// S returns entities' position.
func (s *Step) S() *Shape {
	return &s.Shape
}

// Init do nothing
func (s *Step) Init(m *Model) {
	s.Base.Init(STEP, m)
}

// From returns where Step comes from
func (s *Step) From() Entity {
	return s.FromNode
}

// To returns where Step goes to
func (s *Step) To() Entity {
	return s.ToNode
}

// Cost is calculated by distance
func (s *Step) Cost() float64 {
	return s.Shape.Dist() / Const.Human.Speed
}

// Resolve set reference from id.
func (s *Step) Resolve(args ...Entity) {
}

// CheckDelete check remaining reference.
func (s *Step) CheckDelete() error {
	return nil
}

// BeforeDelete delete selt from related Locationable.
func (s *Step) BeforeDelete() {
	delete(s.FromNode.OutSteps(), s.ID)
	delete(s.ToNode.InSteps(), s.ID)
}

// Delete removes this entity with related ones.
func (s *Step) Delete() {
	s.M.Delete(s)
}

// String represents status
func (s *Step) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,c=%.2f", s.Type().Short(),
		s.ID, s.FromNode, s.ToNode, s.Cost())
}
