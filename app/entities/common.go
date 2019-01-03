package entities

import (
	"time"

	"github.com/revel/revel"
)

// Base based on gorm.Model
type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `                   json:"-"`
	UpdatedAt time.Time  `                   json:"-"`
	DeletedAt *time.Time `gorm:"index"       json:"-"`
}

// NewBase create new Base
func NewBase(id uint) Base {
	return Base{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Owner means this faciliites in under the control by Player.
type Owner struct {
	Own     *Player `gorm:"-"        json:"-"`
	OwnerID uint    `gorm:"not null" json:"oid"`
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
	if o != nil {
		o.OwnerID = o.Own.ID
	}
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
