package entities

import (
	"fmt"
	"math"
	"math/rand"
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
	Persistence
	Shape

	Level       PlayerType `gorm:"not null"       json:"lv"`
	DisplayName string     `gorm:"not null"       json:"name"`
	LoginID     string     `gorm:"not null;index" json:"-"`
	Password    string     `gorm:"not null"       json:"-"`
	ReRouting   bool       `gorm:"-"              json:"-"`
	Color       int        `gorm:"not null"       json:"color"`

	RailNodes map[uint]*RailNode `gorm:"-" json:"-"`
	RailEdges map[uint]*RailEdge `gorm:"-" json:"-"`
	Stations  map[uint]*Station  `gorm:"-" json:"-"`
	Gates     map[uint]*Gate     `gorm:"-" json:"-"`
	Platforms map[uint]*Platform `gorm:"-" json:"-"`
	RailLines map[uint]*RailLine `gorm:"-" json:"-"`
	LineTasks map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`
	Tracks    map[uint]*Track    `gorm:"-" json:"-"`
}

// NewPlayer create instance
func (m *Model) NewPlayer() *Player {
	o := &Player{
		Base:        m.NewBase(PLAYER),
		Persistence: NewPersistence(),
		Shape:       NewShapeGroup(),
	}
	h := rand.Float64() * 360
	var r, g, b int
	h0 := int(math.Floor(float64(h) / 60))

	switch h0 {
	case 0:
		r, g, b = 0xFF, int(h/60*0x100), 0x00
	case 1:
		r, g, b = int((h-120)/60*0x100), 0xFF, 0x00
	case 2:
		r, g, b = 0x00, 0xFF, int((h-120)/60*0x100)
	case 3:
		r, g, b = 0x00, int((240-h)/60*0x100), 0xFF
	case 4:
		r, g, b = int((240-h)/60*0x100), 0x00, 0xFF
	case 5:
		r, g, b = 0xFF, 0x00, int((360-h)/60*0x100)
	}

	o.Color = r<<16 + g<<8 + b

	o.O = o
	o.OwnerID = o.ID
	o.Init(m)
	o.Marshal()
	m.Add(o)
	return o
}

// B returns base information of this elements.
func (o *Player) B() *Base {
	return &o.Base
}

// P returns time information for database.
func (o *Player) P() *Persistence {
	return &o.Persistence
}

// S returns entities' position.
func (o *Player) S() *Shape {
	return &o.Shape
}

// ClearTracks eraces track infomation.
func (o *Player) ClearTracks() {
	for _, tr := range o.Tracks {
		tr.Delete()
	}
}

// Init do nothing
func (o *Player) Init(m *Model) {
	o.Base.Init(PLAYER, m)
	o.Shape.Children = []*Shape{}
	o.RailNodes = make(map[uint]*RailNode)
	o.RailEdges = make(map[uint]*RailEdge)
	o.Stations = make(map[uint]*Station)
	o.Gates = make(map[uint]*Gate)
	o.Platforms = make(map[uint]*Platform)
	o.RailLines = make(map[uint]*RailLine)
	o.LineTasks = make(map[uint]*LineTask)
	o.Trains = make(map[uint]*Train)
	o.Tracks = make(map[uint]*Track)
}

// Resolve set reference.
func (o *Player) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			o.RailNodes[obj.ID] = obj
		case *RailEdge:
			o.RailEdges[obj.ID] = obj
		case *Station:
			o.Stations[obj.ID] = obj
		case *Gate:
			o.Gates[obj.ID] = obj
		case *Platform:
			o.Platforms[obj.ID] = obj
		case *RailLine:
			o.RailLines[obj.ID] = obj
		case *LineTask:
			o.LineTasks[obj.ID] = obj
		case *Train:
			o.Trains[obj.ID] = obj
		case *Track:
			o.Tracks[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type %v %+v", obj, obj))
		}
		o.Shape.Append(raw.S())
	}
}

// UnResolve unregisters specified refernce.
func (o *Player) UnResolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			delete(o.RailNodes, obj.ID)
		case *RailEdge:
			delete(o.RailEdges, obj.ID)
		case *Station:
			delete(o.Stations, obj.ID)
		case *Gate:
			delete(o.Gates, obj.ID)
		case *Platform:
			delete(o.Platforms, obj.ID)
		case *RailLine:
			delete(o.RailLines, obj.ID)
		case *LineTask:
			delete(o.LineTasks, obj.ID)
		case *Train:
			delete(o.Trains, obj.ID)
		case *Track:
			delete(o.Tracks, obj.ID)
		default:
			panic(fmt.Errorf("invalid type %v %+v", obj, obj))
		}
		o.Shape.Delete(raw.S())
	}
}

// Marshal do nothing for implementing Resolvable
func (o *Player) Marshal() {
	// do-nothing
}

// UnMarshal set reference from id.
func (o *Player) UnMarshal() {

}

// CheckDelete check remain relation.
func (o *Player) CheckDelete() error {
	return nil
}

// BeforeDelete deletes related reference
func (o *Player) BeforeDelete() {
}

// Delete removes this entity with related ones.
func (o *Player) Delete() {
	o.M.Delete(o)
}

// String represents status
func (o *Player) String() string {
	o.Marshal()
	return fmt.Sprintf("%s(%d):nm=%s,lv=%v:%s", o.Type().Short(),
		o.ID, o.LoginID, o.Level, o.DisplayName)
}

// Short returns short description
func (o *Player) Short() string {
	return fmt.Sprintf("%s(%d)", o.LoginID, o.ID)
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
