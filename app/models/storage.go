package models

import (
	"sync"
)

type staticModel struct {
}

type agentModel struct {
}

type routeTemplate struct {
}

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
