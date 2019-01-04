package entities

import (
	"fmt"
	"time"
)

// RailLine represents how Train should run.
type RailLine struct {
	Base
	Owner

	Name   string             `         json:"name"`
	Tasks  map[uint]*LineTask `gorm:"-" json:"-"`
	Trains map[uint]*Train    `gorm:"-" json:"-"`
}

// NewRailLine create instance
func NewRailLine(id uint, o *Player) *RailLine {
	l := &RailLine{
		Base:  NewBase(id),
		Owner: NewOwner(o),
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
	return LINE
}

// Init makes map
func (l *RailLine) Init() {
	l.Tasks = make(map[uint]*LineTask)
	l.Trains = make(map[uint]*Train)
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
func (l *RailLine) IsIn(center *Point, scale float64) bool {
	for _, lt := range l.Tasks {
		if lt.IsIn(center, scale) {
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
	return l.Base.IsChanged(after)
}

// Reset set status as not changed
func (l *RailLine) Reset() {
	l.Base.Reset()
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
