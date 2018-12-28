package models

import (
	"sync"

	"github.com/yasshi2525/RushHour/app/models/entities"
)

type config struct {
	Residence residence
	Company   company
}

type residence struct {
	Interval  float64
	Capacity  uint
	Randomize float64
}

type company struct {
	Scale float64
}

type staticModel struct {
	Residences []*entities.Residence
	Companies  []*entities.Company
}

type agentModel struct {
}

type routeTemplate struct {
}

var Config config

var StaticModel staticModel
var AgentModel agentModel
var RouteTemplate routeTemplate

var MuStatic sync.RWMutex
var MuAgent sync.RWMutex
var MuRoute sync.Mutex

func InitStorage() {
	StaticModel = staticModel{}
	AgentModel = agentModel{}
	RouteTemplate = routeTemplate{}

	MuStatic = sync.RWMutex{}
	MuAgent = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
