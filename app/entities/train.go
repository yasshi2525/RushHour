package entities

import (
	"fmt"
	"time"
)

// Train carries Human from Station to Station.
type Train struct {
	Base
	Owner

	Capacity uint `gorm:"not null" json:"capacity"`
	// Mobility represents how many Human can get off at the same time.
	Mobility uint    `gorm:"not null" json:"mobility"`
	Speed    float64 `gorm:"not null" json:"speed"`
	Name     string  `gorm:"not null" json:"name"`
	Progress float64 `gorm:"not null" json:"progress"`
	Occupied uint    `gorm:"-"        json:"occupied"`

	task       *LineTask
	Passengers map[uint]*Human `gorm:"-" json:"-"`
	OnRailEdge *RailEdge       `gorm:"-" json:"-"`
	OnPlatform *Platform       `gorm:"-" json:"-"`

	TaskID     uint `         json:"ltid,omitempty"`
	RailEdgeID uint `gorm:"-" json:"reid,omitempty"`
	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewTrain creates instance
func (m *Model) NewTrain(o *Player) *Train {
	t := &Train{
		Base:  NewBase(m.GenID(TRAIN)),
		Owner: NewOwner(o),
	}
	t.Init()
	t.ResolveRef()
	m.Add(t)
	return t
}

// Idx returns unique id field.
func (t *Train) Idx() uint {
	return t.ID
}

// Type returns type of entitiy
func (t *Train) Type() ModelType {
	return TRAIN
}

// Init makes map
func (t *Train) Init() {
	t.Passengers = make(map[uint]*Human)
}

// Pos returns location
func (t *Train) Pos() *Point {
	if t.task == nil {
		return nil
	}
	return t.task.FromNode().Pos().Div(t, t.Progress)
}

// IsIn returns it should be view or not.
func (t *Train) IsIn(x float64, y float64, scale float64) bool {
	return t.Pos().IsIn(x, y, scale)
}

// Resolve set ID from reference
func (t *Train) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			t.Owner = NewOwner(obj)
			obj.Resolve(t)
		case *LineTask:
			t.task = obj
			obj.Resolve(t)
		case *RailEdge:
			t.OnRailEdge = obj
		case *Platform:
			t.OnPlatform = obj
		case *Human:
			t.Passengers[obj.ID] = obj
			t.Occupied++
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	t.ResolveRef()
}

// ResolveRef set id from reference
func (t *Train) ResolveRef() {
	if t.task != nil {
		t.TaskID = t.task.ID
	}
	if t.OnRailEdge != nil {
		t.RailEdgeID = t.OnRailEdge.ID
	}
	if t.OnPlatform != nil {
		t.PlatformID = t.OnPlatform.ID
	}
}

func (t *Train) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Platform:
			t.OnPlatform = nil
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// CheckRemove check remain relation.
func (t *Train) CheckRemove() error {
	return nil
}

// Permits represents Player is permitted to control
func (t *Train) Permits(o *Player) bool {
	return t.Owner.Permits(o)
}

// Task return next field
func (t *Train) Task() *LineTask {
	return t.task
}

// SetTask changes self changed status for backup
func (t *Train) SetTask(v *LineTask) {
	t.task = v
	t.Change()
	v.Resolve(t)
}

func (t *Train) IsNew() bool {
	return t.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (t *Train) IsChanged(after ...time.Time) bool {
	return t.Base.IsChanged(after...)
}

// Reset set status as not changed
func (t *Train) Reset() {
	t.Base.Reset()
}

// String represents status
func (t *Train) String() string {
	t.ResolveRef()
	ostr := ""
	if t.Own != nil {
		ostr = fmt.Sprintf(":%s", t.Own.Short())
	}
	ltstr := ""
	if t.task != nil {
		ltstr = fmt.Sprintf(",lt=%d", t.task.ID)
	}
	posstr := ""
	if t.Pos() != nil {
		posstr = fmt.Sprintf(":%s", t.Pos())
	}
	return fmt.Sprintf("%s(%v):h=%d/%d%s,%%=%.2f%s%s:%s", t.Type().Short(),
		t.ID, len(t.Passengers), t.Capacity, ltstr, t.Progress, posstr, ostr, t.Name)
}
