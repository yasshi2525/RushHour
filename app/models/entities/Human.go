package entities

import "github.com/jinzhu/gorm"

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
	FromRefer uint
	ToRefer   uint
	From      Residence `gorm:"foreignkey:FromRefer"`
	To        Company   `gorm:"foreignkey:ToRefer"`
	On        Standing  `gorm:"type:int"`
}
