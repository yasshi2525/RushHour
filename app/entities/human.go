package entities

import (
	"fmt"
	"math"
)

// Standing is for judgement Human placement on same X, Y
type Standing uint

const (
	// OnGround represents Human still not arrive at Station or get off Train forcefully
	OnGround Standing = iota + 1
	// OnPlatform represents Human enter Station and wait for Train
	OnPlatform
	// OnTrain represents Human ride on Train
	OnTrain
)

// Human commute from Residence to Company by Train
type Human struct {
	Base
	Persistence
	Shape
	Point

	// Avaialble represents how many seconds Human is able to use for moving or staying.
	Available float64 `gorm:"not null" json:"available"`

	// Mobilty represents how many meters Human can move in a second.
	Mobility float64 `gorm:"not null" json:"mobility"`

	// Angle represents where Human looks for. The unit is radian.
	Angle float64 `gorm:"not null" json:"angle"`

	// Lifespan represents how many seconds Human can live.
	// Human will die after specific term with keeping stay
	// in order to save memory resources.
	Lifespan float64 `gorm:"not null" json:"lifespan"`

	// Progress is [0,1] value representing how much Human proceed current task.
	Progress float64 `gorm:"not null" json:"progress"`

	On Standing `gorm:"-" json:"-"`

	Current *Step `gorm:"-" json:"-"`

	From       *Residence `gorm:"-" json:"-"`
	To         *Company   `gorm:"-" json:"-"`
	onPlatform *Platform
	onTrain    *Train
	out        map[uint]*Step

	FromID     uint `gorm:"not null" json:"rid"`
	ToID       uint `gorm:"not null" json:"cid"`
	PlatformID uint `                json:"pid,omitempty"`
	TrainID    uint `                json:"tid,omitempty"`
}

// NewHuman create instance
func (m *Model) NewHuman(o *Player, x float64, y float64) *Human {
	pos := NewPoint(x, y)
	h := &Human{
		Base:        m.NewBase(HUMAN, o),
		Persistence: NewPersistence(),
		Shape:       NewShapeNode(&pos),
		Point:       pos,
	}
	h.Init(m)
	h.Resolve()
	h.Marshal()
	m.Add(h)

	h.GenOutSteps()
	return h
}

// GenOutSteps Generate Step for Human.
// It's depend on where Human stay.
func (h *Human) GenOutSteps() {
	switch h.On {
	case OnGround:
		// h - C for destination
		h.M.NewStep(h, h.To)
		// h -> G
		for _, g := range h.M.Gates {
			h.M.NewStep(h, g)
		}
	case OnPlatform:
		// h - G, P on Human
		h.M.NewStep(h, h.onPlatform)
		h.M.NewStep(h, h.onPlatform.WithGate)
	case OnTrain:
		// do-nothing
	default:
		panic(fmt.Errorf("invalid type: %T %+v", h.On, h.On))
	}
}

// B returns base information of this elements.
func (h *Human) B() *Base {
	return &h.Base
}

// P returns time information for database.
func (h *Human) P() *Persistence {
	return &h.Persistence
}

// S returns entities' position.
func (h *Human) S() *Shape {
	return &h.Shape
}

// Init creates map.
func (h *Human) Init(m *Model) {
	h.Base.Init(HUMAN, m)
	h.M = m
	h.out = make(map[uint]*Step)
}

// OutSteps returns where it can go to
func (h *Human) OutSteps() map[uint]*Step {
	return h.out
}

// InSteps returns where it comes from
func (h *Human) InSteps() map[uint]*Step {
	return nil
}

// Marshal set id from reference
func (h *Human) Marshal() {
	h.FromID = h.From.ID
	h.ToID = h.To.ID
	if h.onPlatform != nil {
		h.PlatformID = h.onPlatform.ID
	}
	if h.onTrain != nil {
		h.TrainID = h.onTrain.ID
	}
}

// UnMarshal set reference from id.
func (h *Human) UnMarshal() {
	h.Resolve(
		h.M.Find(RESIDENCE, h.FromID),
		h.M.Find(COMPANY, h.ToID))
	// nullable fields
	if h.PlatformID != ZERO {
		h.Resolve(h.M.Find(PLATFORM, h.PlatformID))
	}
	if h.TrainID != ZERO {
		h.Resolve(h.M.Find(TRAIN, h.TrainID))
	}
}

// Resolve set reference
func (h *Human) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Residence:
			h.From = obj
			obj.Resolve(h)
		case *Company:
			h.To = obj
			obj.Resolve(h)
		case *Platform:
			h.onPlatform = obj
			obj.Resolve(h)
		case *Train:
			h.onTrain = obj
			obj.Resolve(h)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	h.Marshal()
}

// BeforeDelete deltes refernce of related entity
func (h *Human) BeforeDelete() {
	delete(h.From.Targets, h.ID)
	delete(h.To.Targets, h.ID)
	if h.onPlatform != nil {
		h.onPlatform.Occupied--
		delete(h.onPlatform.Passengers, h.ID)
	}
	if h.onTrain != nil {
		h.onTrain.Occupied--
		delete(h.onTrain.Passengers, h.ID)
	}
}

// UnResolve unregisters specified refernce.
func (h *Human) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	h.Marshal()
}

// CheckDelete check remain relation.
func (h *Human) CheckDelete() error {
	return nil
}

// Delete removes this entity with related ones.
func (h *Human) Delete(force bool) {
	h.M.Delete(h)
}

// TurnTo make Human turn head to dest.
func (h *Human) turnTo(dest Entity) *Human {
	h.Angle = math.Atan2(dest.S().Pos().Y-h.Y, dest.S().Pos().X-h.X)
	return h
}

// Move make Human walk forward to dist.
// If Human exhaust Available time, then stop.
func (h *Human) move(dist float64) *Human {
	capacity := h.Available * h.Mobility

	if dist > capacity {
		dist = capacity
		h.Available = 0
	} else {
		h.Available -= dist / h.Mobility
	}

	h.X += dist * math.Cos(h.Angle)
	h.Y += dist * math.Sin(h.Angle)
	return h
}

// WalkTo make Human walk to dest point.
// If Human cannot reach it, proceed forward as possible.
func (h *Human) WalkTo(dest Entity) *Human {
	h.turnTo(dest).move(h.S().Pos().Dist(dest.S().Pos()))
	return h
}

func (h *Human) GetIn(t *Train) *Human {
	//TODO
	return h
}

func (h *Human) GetOffForce() *Human {
	//TODO
	return h
}

func (h *Human) GetOff(platform *Platform) *Human {
	//TODO
	return h
}

func (h *Human) Enter(from *Gate, to *Platform) *Human {
	//TODO
	return h
}

func (h *Human) Exit(from *Platform, to *Gate) *Human {
	//TODO
	return h
}

func (h *Human) ExitForce() *Human {
	//TODO
	return h
}

func (h *Human) ShouldGetIn(to *Train) bool {
	// TODO
	return false
}

func (h *Human) ShouldGetOff(from *Train) bool {
	// TODO
	return false
}

// OnPlatform return platform human stands
func (h *Human) OnPlatform() *Platform {
	return h.onPlatform
}

// SetOnPlatform changes self changed status for backup
func (h *Human) SetOnPlatform(v *Platform) {
	h.onPlatform = v
	v.Resolve(h)
	h.Change()
}

// OnTrain return next field
func (h *Human) OnTrain() *Train {
	return h.onTrain
}

// SetOnTrain changes self changed status for backup
func (h *Human) SetOnTrain(v *Train) {
	h.onTrain = v
	v.Resolve(h)
	h.Change()
}

// String represents status
func (h *Human) String() string {
	h.Marshal()
	pstr, tstr := "", ""
	if h.onPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", h.onPlatform.ID)
	}
	if h.onTrain != nil {
		tstr = fmt.Sprintf(",t=%d", h.onTrain.ID)
	}

	return fmt.Sprintf("%s(%d):r=%d,c=%d%s%s,a=%.1f,l=%.1f,%%=%.2f:%v",
		h.Type().Short(),
		h.ID, h.FromID, h.ToID, pstr, tstr, h.Available, h.Lifespan, h.Progress, h.Pos())
}
