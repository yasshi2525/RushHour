package entities

import "github.com/jinzhu/gorm"

// Owner means this faciliites in under the control by Player.
type Owner struct {
	OwnerRefer uint
	Owner      Player `gorm:"foreignKey:OwnerRefer"`
}

// RailNode represents rail track as point.
// Station only stands on RailNode
type RailNode struct {
	gorm.Model
	Owner
	Point

	In  []RailEdge
	Out []RailEdge
}

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	gorm.Model
	Owner

	FromRefer uint
	ToRefer   uint
	From      RailNode `gorm:"foreignKey:FromRefer"`
	To        RailNode `gorm:"foreignKey:ToRefer"`
}

// Platform is the base Human wait for Train.
// Platform can enter only through Gate
type Platform struct {
	gorm.Model
	Owner

	On RailNode

	Capacity uint
	Occupied uint
}

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	gorm.Model
	Owner
	Point

	// Num represents how many Human can pass at the same time
	Num uint
	// Mobility represents time one Human pass Gate
	Mobility float64
	// Occupied represents how many Gate are used by Human
	Occupied uint
}

const GateProdist = 10

// Station compose on Platform and Gate
type Station struct {
	gorm.Model
	Owner

	PlatformRefer uint
	GaterRefer    uint
	Platform      `gorm:"foreignKey:PlatformRefer"`
	Gate          `gorm:"foreignKey:GateRefer"`
}
