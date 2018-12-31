package entities

import (
	"time"
)

// Model based on gorm.Model
type Model struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// NewModel create new Model
func NewModel(id uint) Model {
	return Model{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Company is the destination of Human
type Company struct {
	Model
	Junction

	Targets []*Human `gorm:"-"`

	// Scale : if Scale is bigger, more Human destinate Company
	Scale float64 `gorm:"not null"`
}

// ResolvRef do nothing (for implements Resolvable)
func (c *Company) ResolveRef() {
	// do-nothing
}

// Residence generate Human in a period
type Residence struct {
	Model
	Junction

	Targets []*Human `gorm:"-"`

	Capacity  uint    `gorm:"not null"`
	Available float64 `gorm:"not null"`
}

// ResolvRef do nothing (for implements Resolvable)
func (r *Residence) ResolveRef() {
	// do-nothing
}
