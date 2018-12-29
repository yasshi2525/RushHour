package entities

import (
	"time"

	"github.com/jinzhu/gorm"
)

// NewModel create new gorm.Model
func NewModel(id uint) gorm.Model {
	return gorm.Model{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Company is the destination of Human
type Company struct {
	gorm.Model
	Junction

	Targets []Human `gorm:"-"`

	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64
}

// Residence generate Human in a period
type Residence struct {
	gorm.Model
	Junction

	Targets []Human `gorm:"-"`

	Capacity  uint
	Available float64
}
