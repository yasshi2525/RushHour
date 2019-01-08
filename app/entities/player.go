package entities

import (
	"fmt"
	"time"
)

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

	RailNodes map[uint]*RailNode `gorm:"-" json:"-"`
	RailEdges map[uint]*RailEdge `gorm:"-" json:"-"`
}

// NewPlayer create instance
func NewPlayer(id uint) *Player {
	p := &Player{
		Base: NewBase(id),
	}
	p.Init()
	p.ResolveRef()
	return p
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
func (o *Player) IsIn(x float64, y float64, scale float64) bool {
	return true
}

// Init do nothing
func (o *Player) Init() {
	o.RailNodes = make(map[uint]*RailNode)
	o.RailEdges = make(map[uint]*RailEdge)
}

func (o *Player) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			o.RailNodes[obj.ID] = obj
		case *RailEdge:
			o.RailEdges[obj.ID] = obj
		}
	}
}

// ResolveRef do nothing for implementing Resolvable
func (o *Player) ResolveRef() {
	// do-nothing
}

// CheckRemove check remain relation.
func (o *Player) CheckRemove() error {
	return nil
}

// String represents status
func (o *Player) String() string {
	o.ResolveRef()
	return fmt.Sprintf("%s(%d):nm=%s,lv=%v:%s", Meta.Attr[o.Type()].Short,
		o.ID, o.LoginID, o.Level, o.DisplayName)
}

// Short returns short description
func (o *Player) Short() string {
	return fmt.Sprintf("%s(%d)", o.LoginID, o.ID)
}

func (o *Player) IsNew() bool {
	return o.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (o *Player) IsChanged(after ...time.Time) bool {
	return o.Base.IsChanged(after...)
}

// Reset set status as not changed
func (o *Player) Reset() {
	o.Base.Reset()
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
