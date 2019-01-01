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

// Owner means this faciliites in under the control by Player.
type Owner struct {
	Own     *Player `gorm:"-" json:"-"`
	OwnerID uint    `gorm:"not null" json:"owner_id"`
}

// NewOwner create Juntion
func NewOwner(o *Player) Owner {
	return Owner{
		Own:     o,
		OwnerID: o.ID,
	}
}

// ResolveRef resolve ownerID from Owner
func (o *Owner) ResolveRef() {
	o.OwnerID = o.Own.ID
}

// Permits always permits to Admin, Owner.
func (o *Owner) Permits(target *Player) bool {
	switch target.Level {
	case Admin:
		return true
	case Normal:
		return o.Own == target
	case Guest:
		return false
	default:
		revel.AppLog.Errorf("invalid type %T: %+v", target.Level, target)
		return false
	}
}
