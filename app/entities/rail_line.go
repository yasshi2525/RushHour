package entities

import (
	"fmt"
	"time"
)

// RailLine represents how Train should run.
type RailLine struct {
	Base
	Owner

	Name      string             `         json:"name"`
	Tasks     map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`
	Platforms map[uint]*Platform `gorm:"-" json:"-"`

	slow float64
}

// NewRailLine create instance
func NewRailLine(id uint, o *Player, s float64) *RailLine {
	l := &RailLine{
		Base:  NewBase(id),
		Owner: NewOwner(o),
		slow:  s,
	}
	l.Init()
	return l
}

// Idx returns unique id field.
func (l *RailLine) Idx() uint {
	return l.ID
}

// Type returns type of entitiy
func (l *RailLine) Type() ModelType {
	return RAILLINE
}

// Init makes map
func (l *RailLine) Init() {
	l.Tasks = make(map[uint]*LineTask)
	l.Trains = make(map[uint]*Train)
	l.Platforms = make(map[uint]*Platform)
}

// Pos returns location
func (l *RailLine) Pos() *Point {
	sumX, sumY, cnt := 0.0, 0.0, 0.0
	for _, lt := range l.Tasks {
		if pos := lt.Pos(); pos != nil {
			sumX += pos.X
			sumY += pos.Y
			cnt++
		}
	}
	if cnt > 0 {
		return &Point{sumX / cnt, sumY / cnt}
	}
	return nil
}

// IsIn return true when any LineTask is in,
func (l *RailLine) IsIn(x float64, y float64, scale float64) bool {
	for _, lt := range l.Tasks {
		if lt.IsIn(x, y, scale) {
			return true
		}
	}
	return false
}

// Resolve set reference
func (l *RailLine) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			l.Owner = NewOwner(obj)
		case *Platform:
			l.Platforms[obj.ID] = obj
		case *LineTask:
			l.Tasks[obj.ID] = obj
		case *Train:
			l.Trains[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	l.ResolveRef()
}

// ResolveRef set if from reference
func (l *RailLine) ResolveRef() {
}

// CheckRemove check remain relation.
func (l *RailLine) CheckRemove() error {
	return nil
}

// Permits represents Player is permitted to control
func (l *RailLine) Permits(o *Player) bool {
	return l.Owner.Permits(o)
}

// IsChanged returns true when it is changed after Backup()
func (l *RailLine) IsChanged(after ...time.Time) bool {
	return l.Base.IsChanged(after...)
}

// Reset set status as not changed
func (l *RailLine) Reset() {
	l.Base.Reset()
}

// Borders returns head and tail of LineTask.
// Head and tail are nil when LineTask loops
// Tail is undirecting LineTask, that is LineTask.Next is nil
// Head is undirected  LineTask because head of chain is what any other doesn't target
func (l *RailLine) Borders() (*LineTask, *LineTask) {
	var head, tail *LineTask
	for _, lt := range l.Tasks {
		if lt.Before() == nil {
			head = lt
		}
		if lt.Next() == nil {
			tail = lt
		}
		if head != nil && tail != nil {
			return head, tail
		}
	}
	// looped
	return nil, nil
}

// IsRing returns whether LineTask is looping or not
func (l *RailLine) IsRing() bool {
	if len(l.Tasks) == 0 {
		return false
	}
	h, t := l.Borders()
	return h == nil && t == nil && h != t
}

func (l *RailLine) CanRing() bool {
	head, tail := l.Borders()

	// ringed loop can't loop any more
	if head == nil && tail == nil {
		return false
	}

	switch tail.TaskType {
	case OnDeparture:
		// allow: dept -> not dept
		return head.TaskType != OnDeparture && tail.Stay == head.Dept
	case OnMoving:
		fallthrough
	case OnPassing:
		// allow: move/pass -> not dept
		return head.TaskType != OnDeparture && tail.Moving.ToNode == head.Moving.FromNode
	case OnStopping:
		// allow: stop -> dept
		return head.TaskType == OnDeparture && tail.Dest == head.Stay
	default:
		panic(fmt.Errorf("invalid type = %v", tail.TaskType))
	}
}

// String represents status
func (l *RailLine) String() string {
	l.ResolveRef()
	ostr := ""
	if l.Own != nil {
		ostr = fmt.Sprintf(":%s", l.Own.Short())
	}
	posstr := ""
	if l.Pos() != nil {
		posstr = fmt.Sprintf(":%s", l.Pos())
	}
	return fmt.Sprintf("%s(%d):lt=%d%s%s:%s", Meta.Attr[l.Type()].Short,
		l.ID, len(l.Tasks), posstr, ostr, l.Name)
}
