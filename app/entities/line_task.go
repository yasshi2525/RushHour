package entities

import (
	"fmt"
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
	Base
	Owner

	RailLine *RailLine    `gorm:"-"        json:"-"`
	TaskType LineTaskType `gorm:"not null" json:"type"`
	Next     *LineTask    `gorm:"-"        json:"-"`
	Stay     *Platform    `gorm:"-"        json:"-"`
	Moving   *RailEdge    `gorm:"-"        json:"-"`
	Dest     *Platform    `gorm:"-"        json:"-"`

	Trains map[uint]*Train `gorm:"-" json:"-"`

	RailLineID uint `gorm:"not null" json:"lid"`
	NextID     uint `                json:"next,omitempty"`
	StayID     uint `                json:"p1id,omitempty"`
	MovingID   uint `                json:"reid,omitempty"`
	DestID     uint `                json:"p2id,omitempty"`
}

// NewLineTask create instance
func NewLineTask(id uint, l *RailLine) *LineTask {
	return &LineTask{
		Base:   NewBase(id),
		Owner:  l.Owner,
		Trains: make(map[uint]*Train),
	}
}

// Idx returns unique id field.
func (lt *LineTask) Idx() uint {
	return lt.ID
}

// Type returns type of entitiy
func (lt *LineTask) Type() ModelType {
	return LINETASK
}

// Init do nothing
func (lt *LineTask) Init() {
	lt.Trains = make(map[uint]*Train)
}

// Pos returns entities' position
func (lt *LineTask) Pos() *Point {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.Pos()
	default:
		return lt.Moving.Pos()
	}
}

// IsIn returns it should be view or not.
func (lt *LineTask) IsIn(center *Point, scale float64) bool {
	switch lt.TaskType {
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
			obj.Resolve(lt)
		case *LineTask:
			lt.Next = obj
		case *Platform:
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay = obj
				lt.Stay.Resolve(obj)
			default:
				lt.Dest = obj
				lt.Dest.Resolve(obj)
			}
		case *RailEdge:
			lt.Moving = obj
		case *Train:
			lt.Trains[obj.ID] = obj
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay.Resolve(obj)
			default:
				lt.Moving.Resolve(obj)
			}
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}

	lt.ResolveRef()
}

// ResolveRef set id from reference
func (lt *LineTask) ResolveRef() {
	lt.Owner.ResolveRef()
	if lt.RailLine != nil {
		lt.RailLineID = lt.RailLine.ID
	}
	if lt.Next != nil {
		lt.NextID = lt.Next.ID
	}
	if lt.Moving != nil {
		lt.MovingID = lt.Moving.ID
	}
	if lt.Stay != nil {
		lt.StayID = lt.Stay.ID
	}
	if lt.Dest != nil {
		lt.DestID = lt.Dest.ID
	}
}

func (lt *LineTask) UnRef() {
	// TODO impl
}

// Permits represents Player is permitted to control
func (lt *LineTask) Permits(o *Player) bool {
	return lt.Owner.Permits(o)
}

// From represents start point
func (lt *LineTask) From() Locationable {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.From()
	}
}

// To represents end point
func (lt *LineTask) To() Locationable {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.To()
	}
}

// Cost represents distance
func (lt *LineTask) Cost() float64 {
	switch lt.TaskType {
	case OnDeparture:
		return 0
	default:
		return lt.Moving.Cost()
	}
}

// String represents status
func (lt *LineTask) String() string {
	next, stay, moving := "", "", ""
	if lt.Next != nil {
		next = fmt.Sprintf(",next=%d", lt.Next.ID)
	}
	if lt.Stay != nil {
		stay = fmt.Sprintf(",p=%d", lt.Stay.ID)
	}
	if lt.Moving != nil {
		moving = fmt.Sprintf(",re=%d", lt.Moving.ID)
	}

	return fmt.Sprintf("%s(%d):%v,l=%d%s%s%s:%v:%s", Meta.Attr[lt.Type()].Short,
		lt.ID, lt.TaskType, lt.RailLine.ID, next, stay, moving, lt.Pos(), lt.RailLine.Name)
}

func (ltt LineTaskType) String() string {
	switch ltt {
	case OnDeparture:
		return "dept"
	case OnMoving:
		return "move"
	case OnStopping:
		return "stop"
	case OnPassing:
		return "pass"
	}
	return "????"
}
