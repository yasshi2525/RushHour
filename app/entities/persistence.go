package entities

import "time"

// PersistentStatus represents wethear entity requires database update.
type PersistentStatus uint

const (
	// DBNew represents that entity doesn't exist in database.
	DBNew PersistentStatus = iota
	// DBMerged represents that entity is synchronized to database.
	DBMerged
	// DBChanged represents that old entity exists in database and it needs update.
	DBChanged
)

type Persistence struct {
	CreatedAt time.Time  `                   json:"-"`
	UpdatedAt time.Time  `                   json:"-"`
	DeletedAt *time.Time `gorm:"index"       json:"-"`
	// Changed represents it need to update database
	DBStatus PersistentStatus `gorm:"-" json:"-"`
}

func NewPersistence() Persistence {
	return Persistence{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DBStatus:  DBNew,
	}
}

// IsNew represents this entity doesn't exist in database.
func (p *Persistence) IsNew() bool {
	return p.DBStatus == DBNew
}

// IsChanged returns true when it is changed after
func (p *Persistence) IsChanged() bool {
	return p.DBStatus != DBMerged
}

// Reset set status as not changed
func (p *Persistence) Reset() {
	p.DBStatus = DBMerged
}

// Change marks changeness.
func (p *Persistence) Change() {
	// keep DBNew
	if p.DBStatus != DBNew {
		p.DBStatus = DBChanged
	}
}
