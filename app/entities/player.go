package entities

import "fmt"

// PlayerType represents authenticate level
type PlayerType uint

// PlayerType represents authenticate level
const (
	Admin PlayerType = iota + 1
	Normal
	Guest
)

// Player represents user information
type Player struct {
	Model

	Level       PlayerType `gorm:"not null" json:"lv"`
	DisplayName string     `gorm:"not null" json:"name"`
	LoginID     string     `gorm:"not null;index" json:"-"`
	Password    string     `gorm:"not null" json:"-"`
}

// NewPlayer create instance
func NewPlayer(id uint) *Player {
	return &Player{
		Model: NewModel(id),
	}
}

// Idx returns unique id field.
func (o *Player) Idx() uint {
	return o.ID
}

// Init do nothing
func (o *Player) Init() {
	o.Model.Init()
}

// ResolveRef do nothing for implementing Resolvable
func (o *Player) ResolveRef() {
	// do-nothing
}

// String represents status
func (o *Player) String() string {
	return fmt.Sprintf("%s(%s):lv=%v:%s", Meta.Static[PLAYER].Short,
		o.LoginID, o.Level, o.DisplayName)
}

func (pt PlayerType) String() string {
	switch pt {
	case Admin:
		return "admin"
	case Normal:
		return "normal"
	case Guest:
		return "guest"
	}
	return "???"
}
