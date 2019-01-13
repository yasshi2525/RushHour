package entities

import (
	"fmt"
)

// Step represents two Relayable is logically connected.
// Step is out of target for persistence because it can derived by other resources.
type Step struct {
	ID       uint
	M        *Model
	FromNode Relayable
	ToNode   Relayable
}

// NewWalk create new instance and relation to Relayable
func (m *Model) NewStep(f Relayable, t Relayable) *Step {
	s := &Step{
		ID:       m.GenID(STEP),
		FromNode: f,
		ToNode:   t,
	}
	s.Init(m)
	f.OutSteps()[s.ID] = s
	t.InSteps()[s.ID] = s
	m.Add(s)
	return s
}

// Idx returns unique id field.
func (s *Step) Idx() uint {
	return s.ID
}

// Type returns type of entitiy
func (s *Step) Type() ModelType {
	return STEP
}

// Init do nothing
func (s *Step) Init(m *Model) {
	s.M = m
}

// Pos returns center
func (s *Step) Pos() *Point {
	return s.FromNode.Pos().Center(s.ToNode)
}

// IsIn return true when from, to, center is in,
func (s *Step) IsIn(x float64, y float64, scale float64) bool {
	return s.FromNode.Pos().IsInLine(s.ToNode, x, y, scale)
}

// From returns where Step comes from
func (s *Step) From() Indexable {
	return s.FromNode
}

// To returns where Step goes to
func (s *Step) To() Indexable {
	return s.ToNode
}

// Cost is calculated by distance
func (s *Step) Cost() float64 {
	return s.FromNode.Pos().Dist(s.ToNode) * Const.Human.Weight
}

// UnRef delete selt from related Locationable.
func (s *Step) UnRef() {
	delete(s.FromNode.OutSteps(), s.ID)
	delete(s.ToNode.InSteps(), s.ID)
}

func (s *Step) Delete() {
	s.M.Delete(s)
}

// String represents status
func (s *Step) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,c=%.2f", s.Type().Short(),
		s.ID, s.FromNode, s.ToNode, s.Cost())
}
