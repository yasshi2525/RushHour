package entities

import "github.com/jinzhu/gorm"

// Owneable means this faciliites in under the control by Player.
type Ownable struct {
	OwnerID uint
	Owner   *Player `gorm:"foreignKey:OwnerID"`
}

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	gorm.Model
	Ownable
	Point

	In  []RailEdge
	Out []RailEdge
}

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	gorm.Model
	Ownable

	FromID uint
	ToID   uint
	From   *RailNode
	To     *RailNode
}

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	gorm.Model
	Ownable
	Junction `gorm:"-"`

	In       *Station
	On       *RailNode
	Capacity uint
	Occupied uint
}

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	gorm.Model
	Ownable
	Junction `gorm:"-"`

	In *Station
	// Num represents how many Human can pass at the same time
	Num uint
	// Mobility represents time one Human pass Gate.
	Mobility float64
	// Occupied represents how many Gate are used by Human.
	Occupied uint
}

// Station compose on Platform and Gate
type Station struct {
	gorm.Model
	Ownable

	Name       string
	PlatformID uint
	GateID     uint
	Platform   *Platform
	Gate       *Gate
}

// LineTaskType represents the state what Train should do now.
type LineTaskType uint

const (
	// OnDeparture represents the state that Train waits for departure in Station.
	OnDeparture LineTaskType = iota
	// OnMoving represents the state that Train runs to next RailNode.
	OnMoving
	// OnStopping represents the state that Train stops to next Station.
	OnStopping
	// OnPassing represents the state that Train passes to next Station.
	OnPassing
)

// LineTask is the element of Line.
// The chain of LineTask represents Line structure.
type LineTask struct {
	gorm.Model
	Ownable

	Type   LineTaskType
	NextID uint
	Next   *LineTask
}

type Line struct {
	gorm.Model
	Ownable

	Name  string
	Tasks []*LineTask
}

// Train carries Human from Station to Station.
type Train struct {
	gorm.Model
	Ownable
	Point
	Capacity uint
	// Mobility represents how many Human can get off at the same time.
	Mobility uint
	Speed    float64
	Name     string
}
