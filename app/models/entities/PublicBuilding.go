package entities

import (
	"github.com/jinzhu/gorm"
)

// Company is the destination of Human
type Company struct {
	gorm.Model
	Junction
	Targets []Human `gorm:"-"`

	// Scale : if Scale is bigger, more Human destinate Company
	Scale uint
}

// Residence generate Human in a period
type Residence struct {
	gorm.Model
	Junction
	Targets []Human `gorm:"-"`

	capacity uint
}
