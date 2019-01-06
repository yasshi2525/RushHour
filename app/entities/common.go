package entities

import (
	"time"

	"github.com/revel/revel"
)

type PersistentStatus uint

const (
	DBNew PersistentStatus = iota
	DBMerged
	DBChanged
)

// Base based on gorm.Model
type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `                   json:"-"`
	UpdatedAt time.Time  `                   json:"-"`
	DeletedAt *time.Time `gorm:"index"       json:"-"`
	// Changed represents it need to update database
	DBStatus PersistentStatus `gorm:"-" json:"-"`
	// ChangedAt represents when model is changed. (UpdateAt is for gorm)
	ChangedAt time.Time `gorm:"-" json:"-"`
}

// NewBase create new Base
func NewBase(id uint) Base {
	return Base{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ChangedAt: time.Now(),
	}
}

func (base *Base) IsNew() bool {
	return base.DBStatus == DBNew
}

// IsChanged returns true when it is changed after
func (base *Base) IsChanged(after ...time.Time) bool {
	if len(after) > 0 {
		return base.ChangedAt.Sub(after[0]) > 0
	}
	return base.DBStatus != DBMerged
}

// Reset set status as not changed
func (base *Base) Reset() {
	base.DBStatus = DBMerged
}

// Change marks changeness.
func (base *Base) Change() {
	// keep DBNew
	if base.DBStatus != DBNew {
		base.DBStatus = DBChanged
	}
	base.ChangedAt = time.Now()
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

// Permits always permits to Admin, Owner.
func (o Owner) Permits(target *Player) bool {
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
