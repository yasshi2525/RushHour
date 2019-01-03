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

	From       *Residence `gorm:"-" json:"-"`
	To         *Company   `gorm:"-" json:"-"`
	OnPlatform *Platform  `gorm:"-" json:"-"`
	OnTrain    *Train     `gorm:"-" json:"-"`
	On         Standing   `gorm:"-" json:"-"`
	out        map[uint]*Step

	FromID     uint `gorm:"not null" json:"rid"`
	ToID       uint `gorm:"not null" json:"cid"`
	PlatformID uint `                json:"pid,omitempty"`
	TrainID    uint `                json:"tid,omitempty"`
}

// NewHuman create instance
func NewHuman(id uint, x float64, y float64) *Human {
	return &Human{
		Base:  NewBase(id),
		Point: NewPoint(x, y),
		out:   make(map[uint]*Step),
	}
}

// Idx returns unique id field.
func (h *Human) Idx() uint {
	return h.ID
}

// Type returns type of entitiy
func (h *Human) Type() ModelType {
	return HUMAN
}

// Init creates map.
func (h *Human) Init() {
	h.out = make(map[uint]*Step)
}

// Pos returns entities' position
func (h *Human) Pos() *Point {
	return &h.Point
}

// IsIn returns it should be view or not.
func (h *Human) IsIn(center *Point, scale float64) bool {
	return h.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (h *Human) Out() map[uint]*Step {
	return h.out
}

// In returns where it comes from
func (h *Human) In() map[uint]*Step {
	return nil
}

// ResolveRef set id from reference
func (h *Human) ResolveRef() {
	h.FromID = h.From.ID
	h.ToID = h.To.ID
	h.PlatformID = h.OnPlatform.ID
	h.TrainID = h.OnTrain.ID
}

// Resolve set reference
func (h *Human) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Residence:
			h.From = obj
			obj.Resolve(h)
		case *Company:
			h.To = obj
			obj.Resolve(h)
		case *Platform:
			h.OnPlatform = obj
			obj.Resolve(h)
		case *Train:
			h.OnTrain = obj
			obj.Resolve(h)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	h.ResolveRef()
}

// UnRef deltes refernce of related entity
func (h *Human) UnRef() {
	delete(h.From.Targets, h.ID)
	delete(h.To.Targets, h.ID)
	if h.OnPlatform != nil {
		h.OnPlatform.Occupied--
		delete(h.OnPlatform.Passenger, h.ID)
	}
	if h.OnTrain != nil {
		h.OnTrain.Occupied--
		delete(h.OnTrain.Passenger, h.ID)
	}
}

// TurnTo make Human turn head to dest.
func (h *Human) turnTo(dest *Point) *Human {
	h.Angle = math.Atan2(dest.Y-h.Y, dest.X-h.X)
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
func (h *Human) WalkTo(dest *Point) *Human {
	h.turnTo(dest).move(h.Dist(dest))
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

// String represents status
func (h *Human) String() string {
	pstr, tstr := "", ""
	if h.OnPlatform != nil {
		pstr = fmt.Sprintf(",p=%d", h.OnPlatform.ID)
	}
	if h.OnTrain != nil {
		tstr = fmt.Sprintf(",t=%d", h.OnTrain.ID)
	}

	return fmt.Sprintf("%s(%d):r=%d,c=%d%s%s,a=%.1f,l=%.1f,%%=%.2f:%v",
		Meta.Attr[h.Type()].Short,
		h.ID, h.From.ID, h.To.ID, pstr, tstr, h.Available, h.Lifespan, h.Progress, h.Pos())
}
