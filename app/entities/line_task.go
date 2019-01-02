package entities

import (
	"github.com/revel/revel"
)

// LineTaskType represents the state what Train should do now.
type LineTaskType uint

const (
	// OnDeparture represents the state that Train waits for departure in Station.
	OnDeparture LineTaskType = iota + 1
	// OnMoving represents the state that Train runs to next RailNode.
	OnMoving
	// OnStopping represents the state that Train stops to next Station.
	OnStopping
	// OnPassing represents the state that Train passes to next Station.
	OnPassing
)

// LineTask is the element of Line.
// The chain of LineTask represents Line structure.
type LineTask struct {
	Model
	Owner

	RailLine *RailLine    `gorm:"-" json:"-"`
	Type     LineTaskType `gorm:"not null"`
	Next     *LineTask    `gorm:"-" json:"-"`

	Stay   *Platform `gorm:"-" json:"-"`
	Moving *RailEdge `gorm:"-" json:"-"`

	RailLineID uint `gorm:"not null"`
	NextID     uint
	StayID     uint
	MovingID   uint
}

// NewLineTask create instance
func NewLineTask(id uint, l *RailLine) *LineTask {
	return &LineTask{
		Model: NewModel(id),
		Owner: l.Owner,
	}
}

// Idx returns unique id field.
func (lt *LineTask) Idx() uint {
	return lt.ID
}

// Init do nothing
func (lt *LineTask) Init() {
	lt.Model.Init()
	lt.Owner.Init()
}

// Pos returns entities' position
func (lt *LineTask) Pos() *Point {
	switch lt.Type {
	case OnDeparture:
		return lt.Stay.Pos()
	default:
		return lt.Moving.Pos()
	}
}

// IsIn returns it should be view or not.
func (lt *LineTask) IsIn(center *Point, scale float64) bool {
	switch lt.Type {
	case OnDeparture:
		return lt.Stay.IsIn(center, scale)
	default:
		return lt.Moving.IsIn(center, scale)
	}
}

// Resolve set reference
func (lt *LineTask) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailLine:
			lt.Owner, lt.RailLine = obj.Owner, obj
		case *LineTask:
			lt.Next = obj
		case *Platform:
			lt.Stay = obj
		case *RailEdge:
			lt.Moving = obj
		default:
			revel.AppLog.Warnf("invalid type %T: %+v", obj, obj)
		}
	}

	lt.ResolveRef()
}

// ResolveRef set id from reference
func (lt *LineTask) ResolveRef() {
	lt.Owner.ResolveRef()
	lt.RailLineID = lt.RailLine.ID
	if lt.Next != nil {
		lt.NextID = lt.Next.ID
	}
	if lt.Moving != nil {
		lt.MovingID = lt.Moving.ID
	}
	if lt.Stay != nil {
		lt.StayID = lt.Stay.ID
	}
}

// Permits represents Player is permitted to control
func (lt *LineTask) Permits(o *Player) bool {
	return lt.Owner.Permits(o)
}

// From represents start point
func (lt *LineTask) From() Locationable {
	switch lt.Type {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.From()
	}
}

// To represents end point
func (lt *LineTask) To() Locationable {
	switch lt.Type {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.To()
	}
}

// Cost represents distance
func (lt *LineTask) Cost() float64 {
	switch lt.Type {
	case OnDeparture:
		return 0
	default:
		return lt.Moving.Cost()
	}
}
