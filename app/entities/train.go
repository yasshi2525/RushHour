package entities

import (
	"fmt"
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
	Passenger  map[uint]*Human `gorm:"-" json:"-"`
	OnRailEdge *RailEdge       `gorm:"-" json:"-"`
	OnPlatform *Platform       `gorm:"-" json:"-"`

	TaskID     uint `json:"ltid,omitempty"`
	RailEdgeID uint `gorm:"-" json:"reid,omitempty"`
	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewTrain creates instance
func NewTrain(id uint, o *Player) *Train {
	return &Train{
		Base:      NewBase(id),
		Owner:     NewOwner(o),
		Passenger: make(map[uint]*Human),
	}
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
	t.Passenger = make(map[uint]*Human)
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
			t.Passenger[obj.ID] = obj
			t.Occupied++
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	t.ResolveRef()
}

// ResolveRef set id from reference
func (t *Train) ResolveRef() {
	t.Owner.ResolveRef()
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

// Permits represents Player is permitted to control
func (t *Train) Permits(o *Player) bool {
	return t.Owner.Permits(o)
}

// String represents status
func (t *Train) String() string {
	ltstr := ""
	if t.Task != nil {
		ltstr = fmt.Sprintf(",lt=%d", t.Task.ID)
	}
	return fmt.Sprintf("%s(%v):h=%d/%d%s,%%=%.2f:%v:%s", Meta.Attr[t.Type()].Short,
		t.ID, len(t.Passenger), t.Capacity, ltstr, t.Progress, t.Pos(), t.Name)
}
