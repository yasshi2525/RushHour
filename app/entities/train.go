package entities

import (
	"fmt"
	"time"
)

// Train carries Human from Station to Station.
type Train struct {
	Base
	Point
	Owner

	Capacity int `gorm:"not null" json:"capacity"`
	// Mobility represents how many Human can get off at the same time.
	Mobility int     `gorm:"not null" json:"mobility"`
	Speed    float64 `gorm:"not null" json:"speed"`
	Name     string  `gorm:"not null" json:"name"`
	Progress float64 `gorm:"not null" json:"progress"`
	Occupied int     `gorm:"-"        json:"occupied"`

	M          *Model    `gorm:"-" json:"-"`
	OnRailEdge *RailEdge `gorm:"-" json:"-"`
	OnPlatform *Platform `gorm:"-" json:"-"`
	task       *LineTask
	Passengers map[uint]*Human `gorm:"-" json:"-"`

	TaskID     uint `         json:"ltid,omitempty"`
	RailEdgeID uint `gorm:"-" json:"reid,omitempty"`
	PlatformID uint `gorm:"-" json:"pid,omitempty"`
}

// NewTrain creates instance
func (m *Model) NewTrain(o *Player, name string) *Train {
	t := &Train{
		Base:     NewBase(m.GenID(TRAIN)),
		Owner:    NewOwner(o),
		Capacity: Const.Train.Capacity,
		Mobility: Const.Train.Mobility,
		Speed:    Const.Train.Speed,
		Name:     name,
	}
	t.Init(m)
	t.Marshal()
	m.Add(t)
	return t
}

func (t *Train) UnLoad() {
	for _, h := range t.Passengers {
		h.Point = *t.Pos().Rand(Const.Train.Randomize)
		h.onTrain = nil
		h.TrainID = ZERO
		t.Occupied--
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
func (t *Train) Init(m *Model) {
	t.M = m
	t.Passengers = make(map[uint]*Human)
}

// Pos returns location
func (t *Train) Pos() *Point {
	if t.task == nil {
		return &Point{}
	}
	return t.task.FromNode().Pos().Div(t, t.Progress)
}

// IsIn returns it should be view or not.
func (t *Train) IsIn(x float64, y float64, scale float64) bool {
	if t.task == nil {
		return false
	}
	return t.Pos().IsIn(x, y, scale)
}

func (t *Train) SetTask(lt *LineTask) {
	if len(t.Passengers) > 0 {
		panic(fmt.Errorf("try to set task to Train with passengers: %v", t))
	}
	t.task = lt
	if lt != nil {
		t.TaskID = lt.ID
	} else {
		t.TaskID = ZERO
	}
	pos := t.Pos()
	t.X, t.Y = pos.X, pos.Y
	t.Change()
	t.Marshal()
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
	t.Marshal()
}

// Marshal set id from reference
func (t *Train) Marshal() {
	if t.task != nil {
		t.TaskID = t.task.ID
	}
	if t.OnRailEdge != nil {
		t.RailEdgeID = t.OnRailEdge.ID
	} else {
		t.RailEdgeID = ZERO
	}
	if t.OnPlatform != nil {
		t.PlatformID = t.OnPlatform.ID
	} else {
		t.PlatformID = ZERO
	}
}

func (t *Train) UnMarshal() {
	t.Resolve(t.M.Find(PLAYER, t.OwnerID))
	// nullable fields
	if t.TaskID != ZERO {
		t.Resolve(t.M.Find(LINETASK, t.TaskID))
	}
}

func (t *Train) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// CheckDelete check remain relation.
func (t *Train) CheckDelete() error {
	return nil
}

func (t *Train) BeforeDelete() {
	t.UnLoad()
}

func (t *Train) Delete() {
	t.M.Delete(t)
}

// Permits represents Player is permitted to control
func (t *Train) Permits(o *Player) bool {
	return t.Owner.Permits(o)
}

// Task return next field
func (t *Train) Task() *LineTask {
	return t.task
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
	t.Marshal()
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
