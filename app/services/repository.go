package services

import (
	"sync"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/route"
)

// Model is contained all data for gaming
var Model *entities.Model

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate *route.Payload

// MuStatic is mutex lock for Static
var MuStatic sync.RWMutex

// MuDynamic is mutex lock for Dynamic
var MuDynamic sync.RWMutex

// MuRoute is mutex lock for routing
var MuRoute sync.Mutex

// InitLock must prepare first.
func InitLock() {
	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}

// InitRepository initialize storage
func InitRepository() {
	entities.InitType()
	Model = entities.NewModel()
}

// GenStep generate Step and resister it
func GenStep(from entities.Relayable, to entities.Relayable) *entities.Step {
	s := Model.NewWalkStep(from, to, Config.Human.Weight)

	return s
}
