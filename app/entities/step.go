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
	cost float64
	By   *LineTask
}

// NewWalkStep create new instance and relation to Relayable
func NewWalkStep(id uint, f Relayable, t Relayable, cost float64) *Step {
	return newStep(id, f, t, cost)
}

// NewTrainStep create new instance and relation to Platform
func NewTrainStep(id uint, lt *LineTask, dept *Platform, dest *Platform, c float64) *Step {
	s := newStep(id, dept, dest, c)
	s.By = lt
	return s
}

func newStep(id uint, f Relayable, t Relayable, c float64) *Step {
	step := &Step{
		ID:   id,
		from: f,
		to:   t,
		cost: c,
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
	return s.cost
}

// UnRef delete selt from related Locationable.
func (s *Step) UnRef() {
	delete(s.from.OutStep(), s.ID)
	delete(s.to.InStep(), s.ID)
}

// String represents status
func (s *Step) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v", Meta.Attr[s.Type()].Short,
		s.ID, s.from, s.to)
}
