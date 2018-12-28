package entities

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Standing is for judgement Human placement on same X, Y
type Standing uint

const (
	// OnGround represents Human still not arrive at Station or get off Train forcefully
	OnGround Standing = iota
	// OnPlatform represents Human enter Station and wait for Train
	OnPlatform
	// OnTrain represents Human ride on Train
	OnTrain
)

// Human commute from Residence to Company by Train
type Human struct {
	gorm.Model
	Point

	// Avaialble represents how many seconds Human is able to use for moving or staying.
	Available float64

	// Mobilty represents how many meters Human can move in a second.
	Mobility float64

	// Angle represents where Human looks for. The unit is radian.
	Angle float64

	// Lifespan represents how many seconds Human can live.
	// Human will die after specific term with keeping stay
	// in order to save memory resources.
	Lifespan float64

	// Progress is [0,1] value representing how much Human proceed current task.
	Progress float64

	FromID     uint
	ToID       uint
	From       *Residence
	To         *Company
	OnPlatform *Platform
	OnTrain    *Train
	On         Standing
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
