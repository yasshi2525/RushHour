package entities

import (
	"fmt"
)

// Step represents two Relayable is logically connected.
// Step is out of target for persistence because it can derived by other resources.
type Step struct {
	ID   uint
	from Relayable
	to   Relayable

	// only for walk step
	weight float64

	// only for train step
	By        *LineTask
	Transport float64
}

// NewWalkStep create new instance and relation to Relayable
func NewWalkStep(id uint, f Relayable, t Relayable, w float64) *Step {
	s := NewStep(id, f, t)
	s.weight = w
	return s
}

func NewStep(id uint, f Relayable, t Relayable) *Step {
	step := &Step{
		ID:   id,
		from: f,
		to:   t,
	}
	step.Init()
	f.OutStep()[step.ID] = step
	t.InStep()[step.ID] = step
	return step
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
func (s *Step) Init() {
}

// Pos returns center
func (s *Step) Pos() *Point {
	return s.from.Pos().Center(s.to)
}

// From returns where Step comes from
func (s *Step) From() Indexable {
	return s.from
}

// To returns where Step goes to
func (s *Step) To() Indexable {
	return s.to
}

// Cost is calculated by distance
func (s *Step) Cost() float64 {
	if s.By != nil {
		return s.Transport
	}
	return s.from.Pos().Dist(s.to) * s.weight
}

// UnRef delete selt from related Locationable.
func (s *Step) UnRef() {
	delete(s.from.OutStep(), s.ID)
	delete(s.to.InStep(), s.ID)
}

// String represents status
func (s *Step) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,c=%.2f", Meta.Attr[s.Type()].Short,
		s.ID, s.from, s.to, s.Cost())
}
