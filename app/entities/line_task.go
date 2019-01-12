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

	RailLine  *RailLine    `gorm:"-"        json:"-"`
	TaskType  LineTaskType `gorm:"not null" json:"type"`
	before    *LineTask
	next      *LineTask
	Stay      *Platform       `gorm:"-" json:"-"`
	Dept      *Platform       `gorm:"-" json:"-"`
	Moving    *RailEdge       `gorm:"-" json:"-"`
	Dest      *Platform       `gorm:"-" json:"-"`
	Trains    map[uint]*Train `gorm:"-" json:"-"`
	OverSteps map[uint]*Step  `gorm:"-" json:"-"`

	RailLineID uint `gorm:"not null" json:"lid"`
	BeforeID   uint `gorm:"-"        json:"before,omitempty"`
	NextID     uint `                json:"next,omitempty"`
	StayID     uint `                json:"psid,omitempty"`
	DeptID     uint `gorm:"-"        json:"p1id,omitempty"`
	MovingID   uint `                json:"reid,omitempty"`
	DestID     uint `gorm:"-"        json:"p2id,omitempty"`
}

// NewLineTaskDept create "dept"
func (m *Model) NewLineTaskDept(l *RailLine, p *Platform, tail ...*LineTask) *LineTask {
	lt := &LineTask{
		Base:     NewBase(m.GenID(LINETASK)),
		TaskType: OnDeparture,
	}
	lt.Init()
	lt.Resolve(l.Own, l, p)
	lt.ResolveRef()
	if len(tail) > 0 && tail[0] != nil {
		tail[0].SetNext(lt)
	}
	m.Add(lt)
	return lt
}

// NewLineTask creates "stop" or "pass" or "moving"
func (m *Model) NewLineTask(tail *LineTask, re *RailEdge, pass bool) *LineTask {
	lt := &LineTask{
		Base:     NewBase(m.GenID(LINETASK)),
	}
	lt.Init()
	if re.ToNode.OverPlatform == nil {
		lt.TaskType = OnMoving
	} else {
		if pass {
			lt.TaskType = OnPassing
		} else {
			lt.TaskType = OnStopping
		}
	}
	lt.Resolve(tail.Own, tail.RailLine, re)
	lt.ResolveRef()
	tail.SetNext(lt)
	m.Add(lt)
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
	lt.OverSteps = make(map[uint]*Step)
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
		case *Player:
			lt.Owner = NewOwner(obj)
			obj.Resolve(lt)
		case *Platform:
			lt.Stay = obj
			lt.Dept = obj
			lt.RailLine.Resolve(obj)
			obj.Resolve(lt)
			obj.OnRailNode.Resolve(lt.RailLine)
			obj.OnRailNode.OutTasks[lt.ID] = lt
			obj.OnRailNode.InTasks[lt.ID] = lt
		case *RailEdge:
			lt.Moving = obj
			lt.RailLine.Resolve(obj)
			obj.Resolve(lt)
			if p := obj.FromNode.OverPlatform; p != nil {
				lt.Dept = p
				p.Resolve(lt)
			}
			if p := obj.ToNode.OverPlatform; p != nil {
				lt.Dest = p
				p.Resolve(lt)
			}
			obj.FromNode.Resolve(lt.RailLine)
			obj.FromNode.OutTasks[lt.ID] = lt
			obj.ToNode.Resolve(lt.RailLine)
			obj.ToNode.InTasks[lt.ID] = lt
		case *RailLine:
			lt.RailLine = obj
			obj.Resolve(lt)
		case *LineTask:
			lt.next = obj
			if obj != nil {
				obj.before = lt
				obj.ResolveRef()
			}
		case *Train:
			lt.Trains[obj.ID] = obj
			lt.RailLine.Resolve(obj)
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay.Resolve(obj)
			default:
				lt.Moving.Resolve(obj)
			}
		case *Step:
			lt.OverSteps[obj.ID] = obj
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
	if lt.before != nil {
		lt.BeforeID = lt.before.ID
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
	if lt.Dept != nil {
		lt.DeptID = lt.Dept.ID
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
func (lt *LineTask) From() Indexable {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.FromNode
	}
}

// To represents end point
func (lt *LineTask) To() Indexable {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay
	default:
		return lt.Moving.ToNode
	}
}

// FromNode represents start point
func (lt *LineTask) FromNode() *RailNode {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.OnRailNode
	default:
		return lt.Moving.FromNode
	}
}

// ToNode represents end point
func (lt *LineTask) ToNode() *RailNode {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.OnRailNode
	default:
		return lt.Moving.ToNode
	}
}

// Cost represents distance
func (lt *LineTask) Cost() float64 {
	switch lt.TaskType {
	case OnDeparture:
		return 0
	default:
		cost := lt.Moving.Cost()
		if lt.before.TaskType == OnDeparture {
			cost /= Const.Train.Slowness
		}
		if lt.TaskType == OnStopping {
			cost /= Const.Train.Slowness
		}
		return cost
	}
}

// Before return before field
func (lt *LineTask) Before() *LineTask {
	return lt.before
}

// Next return next field
func (lt *LineTask) Next() *LineTask {
	return lt.next
}

// SetNext changes self changed status for backup
func (lt *LineTask) SetNext(v *LineTask) {
	lt.next = v
	if v != nil {
		v.before = lt
	}
	lt.Change()
}

func (lt *LineTask) IsNew() bool {
	return lt.Base.IsNew()
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
	before, next, stay, dept, moving, dest := "", "", "", "", "", ""
	if lt.before != nil {
		before = fmt.Sprintf(",before=%d", lt.before.ID)
	}
	if lt.next != nil {
		next = fmt.Sprintf(",next=%d", lt.next.ID)
	}
	if lt.Stay != nil {
		stay = fmt.Sprintf(",stay=%d", lt.Stay.ID)
	}
	if lt.Dept != nil {
		dept = fmt.Sprintf(",dept=%d", lt.Dept.ID)
	}
	if lt.Moving != nil {
		moving = fmt.Sprintf(",re=%d", lt.Moving.ID)
	}
	if lt.Dest != nil {
		dest = fmt.Sprintf(",dest=%d", lt.Dest.ID)
	}
	posstr := ""
	if lt.Pos() != nil {
		posstr = fmt.Sprintf(":%s", lt.Pos())
	}
	nmstr := ""
	if lt.RailLine != nil {
		nmstr = fmt.Sprintf(":%s", lt.RailLine.Name)
	}
	return fmt.Sprintf("%s(%d):%v,l=%d%s%s%s%s%s%s%s%s%s", lt.Type().Short(),
		lt.ID, lt.TaskType, lt.RailLineID, before, next, stay, dept, moving, dest,
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
