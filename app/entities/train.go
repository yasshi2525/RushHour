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

	Task       *LineTask       `gorm:"-" json:"-"`
	Passengers map[uint]*Human `gorm:"-" json:"-"`
	OnRailEdge *RailEdge       `gorm:"-" json:"-"`
	OnPlatform *Platform       `gorm:"-" json:"-"`

	TaskID     uint `json:"ltid,omitempty"`
	RailEdgeID uint `gorm:"-" json:"reid,omitempty"`
	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewTrain creates instance
func NewTrain(id uint, o *Player) *Train {
	t := &Train{
		Base:  NewBase(id),
		Owner: NewOwner(o),
	}
	t.Init()
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
	if t.Task == nil {
		return nil
	}
	from, to := t.Task.From().Pos(), t.Task.To().Pos()
	return from.Div(to, t.Progress)
}

// IsIn returns it should be view or not.
func (t *Train) IsIn(center *Point, scale float64) bool {
	return t.Pos().IsIn(center, scale)
}

// Resolve set ID from reference
func (t *Train) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *LineTask:
			t.Owner, t.Task = obj.Owner, obj
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
	if t.Task != nil {
		t.TaskID = t.Task.ID
	}
	if t.OnRailEdge != nil {
		t.RailEdgeID = t.OnRailEdge.ID
	}
	if t.OnPlatform != nil {
		t.PlatformID = t.OnPlatform.ID
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

// IsChanged returns true when it is changed after Backup()
func (t *Train) IsChanged(after ...time.Time) bool {
	return t.Base.IsChanged(after)
}

// Reset set status as not changed
func (t *Train) Reset() {
	t.Base.Reset()
}

// String represents status
func (t *Train) String() string {
	ostr := ""
	if t.Own != nil {
		ostr = fmt.Sprintf(":%s", t.Own.Short())
	}
	ltstr := ""
	if t.Task != nil {
		ltstr = fmt.Sprintf(",lt=%d", t.Task.ID)
	}
	posstr := ""
	if t.Pos() != nil {
		posstr = fmt.Sprintf(":%s", t.Pos())
	}
	return fmt.Sprintf("%s(%v):h=%d/%d%s,%%=%.2f%s%s:%s", Meta.Attr[t.Type()].Short,
		t.ID, len(t.Passengers), t.Capacity, ltstr, t.Progress, posstr, ostr, t.Name)
}
