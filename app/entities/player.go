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
	Base

	Level       PlayerType `gorm:"not null"       json:"lv"`
	DisplayName string     `gorm:"not null"       json:"name"`
	LoginID     string     `gorm:"not null;index" json:"-"`
	Password    string     `gorm:"not null"       json:"-"`
}

// NewPlayer create instance
func NewPlayer(id uint) *Player {
	return &Player{
		Base: NewBase(id),
	}
}

// Idx returns unique id field.
func (o *Player) Idx() uint {
	return o.ID
}

// Type returns type of entitiy
func (o *Player) Type() ModelType {
	return PLAYER
}

// Pos returns nil
func (o *Player) Pos() *Point {
	return nil
}

// IsIn always returns true in order for user to view other Player
func (o *Player) IsIn(center *Point, scale float64) bool {
	return true
}

// Init do nothing
func (o *Player) Init() {
}

// ResolveRef do nothing for implementing Resolvable
func (o *Player) ResolveRef() {
	// do-nothing
}

// String represents status
func (o *Player) String() string {
	return fmt.Sprintf("%s(%s):lv=%v:%s", Meta.Attr[o.Type()].Short,
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
