package entities

import (
	"time"

	"github.com/revel/revel"
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

// Resolvable set some_id fields from reference.
// Resolvable is for database migration
type Resolvable interface {
	// ResolveRef set id from object
	ResolveRef()
}

// OwnableEntity works as auth level
type OwnableEntity interface {
	// Permits represents Player is permitted to control
	Permits(*Player) bool
}

// Ownable means this faciliites in under the control by Player.
type Ownable struct {
	Owner *Player `gorm:"-" json:"-"`

	OwnerID uint `gorm:"not null" json:"owner_id"`
}

// NewOwnable create Juntion
func NewOwnable(o *Player) Ownable {
	return Ownable{
		Owner:   o,
		OwnerID: o.ID,
	}
}

// ResolveRef resolve ownerID from Owner
func (o *Ownable) ResolveRef() {
	o.OwnerID = o.Owner.ID
}

// Permits always permits to Admin, Owner.
func (o *Ownable) Permits(target *Player) bool {
	switch target.Level {
	case Admin:
		return true
	case Normal:
		return o.Owner == target
	case Guest:
		return false
	default:
		revel.AppLog.Errorf("invalid type %T: %+v", target.Level, target)
		return false
	}
}
