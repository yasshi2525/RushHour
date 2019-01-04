package entities

import (
	"fmt"
	"time"
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
	next     *LineTask    `gorm:"-"        json:"-"`
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

// NewLineTaskDept create "dept"
func NewLineTaskDept(id uint, l *RailLine, p *Platform, tail ...*LineTask) *LineTask {
	lt := &LineTask{
		Base:     NewBase(id),
		Owner:    l.Owner,
		RailLine: l,
		TaskType: OnDeparture,
		Stay:     p,
	}
	lt.Init()
	lt.ResolveRef()
	l.Resolve(lt)
	p.Resolve(lt)
	if len(tail) > 0 {
		tail[0].Resolve(lt)
	}
	return lt
}

// NewLineTask creates "stop" or "pass" or "moving"
func NewLineTask(id uint, tail *LineTask, re *RailEdge, pass ...bool) *LineTask {
	lt := &LineTask{
		Base:     NewBase(id),
		Owner:    tail.Owner,
		RailLine: tail.RailLine,
		Moving:   re,
		Dest:     re.ToNode.OverPlatform,
	}
	if re.ToNode.OverPlatform == nil {
		lt.TaskType = OnMoving
	} else {
		if len(pass) > 0 && pass[0] {
			lt.TaskType = OnPassing
		} else {
			lt.TaskType = OnStopping
		}
	}
	lt.Init()
	lt.ResolveRef()
	tail.RailLine.Resolve(lt)
	re.Resolve(lt)
	if re.ToNode.OverPlatform != nil {
		re.ToNode.OverPlatform.Resolve(lt)
	}
	tail.Resolve(lt)
	return lt
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
		if lt.Stay == nil {
			return nil
		}
		return lt.Stay.Pos()
	default:
		if lt.Moving == nil {
			return nil
		}
		return lt.Moving.Pos()
	}
}

// IsIn returns it should be view or not.
func (lt *LineTask) IsIn(x float64, y float64, scale float64) bool {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.IsIn(x, y, scale)
	default:
		return lt.Moving.IsIn(x, y, scale)
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
			lt.next = obj
		case *Platform:
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay = obj
				obj.Resolve(lt)
				lt.RailLine.Resolve(obj)
			default:
				lt.Dest = obj
				obj.Resolve(lt)
			}
		case *RailEdge:
			lt.Moving = obj
			obj.Resolve(lt)
		case *Train:
			lt.Trains[obj.ID] = obj
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay.Resolve(obj)
			default:
				lt.Moving.Resolve(obj)
			}
			lt.RailLine.Resolve(obj)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}

	lt.ResolveRef()
}

// ResolveRef set id from reference
func (lt *LineTask) ResolveRef() {
	if lt.RailLine != nil {
		lt.RailLineID = lt.RailLine.ID
	}
	if lt.next != nil {
		lt.NextID = lt.next.ID
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

// UnRef remove related refernce
func (lt *LineTask) UnRef() {
	// TODO impl
}

// CheckRemove check remain relation.
func (lt *LineTask) CheckRemove() error {
	return nil
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

// Next return next field
func (lt *LineTask) Next() *LineTask {
	return lt.next
}

// SetNext changes self changed status for backup
func (lt *LineTask) SetNext(v *LineTask) {
	lt.next = v
	lt.Change()
}

// IsChanged returns true when it is changed after Backup()
func (lt *LineTask) IsChanged(after ...time.Time) bool {
	return lt.Base.IsChanged(after...)
}

// Reset set status as not changed
func (lt *LineTask) Reset() {
	lt.Base.Reset()
}

// String represents status
func (lt *LineTask) String() string {
	lt.ResolveRef()
	ostr := ""
	if lt.Own != nil {
		ostr = fmt.Sprintf(":%s", lt.Own.Short())
	}
	next, stay, moving := "", "", ""
	if lt.next != nil {
		next = fmt.Sprintf(",next=%d", lt.next.ID)
	}
	if lt.Stay != nil {
		stay = fmt.Sprintf(",p=%d", lt.Stay.ID)
	}
	if lt.Moving != nil {
		moving = fmt.Sprintf(",re=%d", lt.Moving.ID)
	}
	posstr := ""
	if lt.Pos() != nil {
		posstr = fmt.Sprintf(":%s", lt.Pos())
	}
	nmstr := ""
	if lt.RailLine != nil {
		nmstr = fmt.Sprintf(":%s", lt.RailLine.Name)
	}
	return fmt.Sprintf("%s(%d):%v,l=%d%s%s%s%s%s%s", Meta.Attr[lt.Type()].Short,
		lt.ID, lt.TaskType, lt.RailLineID, next, stay, moving,
		posstr, ostr, nmstr)
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
