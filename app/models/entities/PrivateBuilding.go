package entities

import "github.com/jinzhu/gorm"

// Owner means this faciliites in under the control by Player.
type Ownable struct {
	OwnerRefer uint
	Owner      *Player `gorm:"foreignKey:OwnerRefer"`
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

	FromRefer uint
	ToRefer   uint
	From      *RailNode `gorm:"foreignKey:FromRefer"`
	To        *RailNode `gorm:"foreignKey:ToRefer"`
}

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	gorm.Model
	Ownable

	On       *RailNode
	Capacity uint
	Occupied uint
}

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	gorm.Model
	Ownable
	Point

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

	Name          string
	PlatformRefer uint
	GaterRefer    uint
	Platform      *Platform `gorm:"foreignKey:PlatformRefer"`
	Gate          *Gate     `gorm:"foreignKey:GateRefer"`
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

	Type      LineTaskType `gorm:"type:int"`
	NextRefer uint
	Next      *LineTask `gorm:"foreignKey:NextRefer"`
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
