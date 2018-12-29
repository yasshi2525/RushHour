package models

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/models/entities"
	validator "gopkg.in/go-playground/validator.v9"
)

type nextID struct {
	Residence uint32
	Company   uint32
}

type config struct {
	Residence residence
	Company   company
}

type residence struct {
	Interval  float64 `validate:"gt=0"`
	Capacity  uint    `validate:"min=1"`
	Randomize float64 `validate:"min=0"`
}

type company struct {
	Scale float64 `validate:"gt=0"`
}

type staticModel struct {
	Residences []entities.Residence
	Companies  []entities.Company
}

type agentModel struct {
}

type routeTemplate struct {
}

var Config config
var NextID nextID

var StaticModel staticModel
var AgentModel agentModel
var RouteTemplate routeTemplate

var MuStatic sync.RWMutex
var MuAgent sync.RWMutex
var MuRoute sync.Mutex

// LoadConf load and validate game.conf
func LoadConf() {
	if _, err := toml.DecodeFile("conf/game.conf", &Config); err != nil {
		revel.AppLog.Errorf("failed to load conf", err)
	}

	if err := validator.New().Struct(Config); err != nil {
		revel.AppLog.Error("%v", err)
	}
}

// InitStorage initialize storage
func InitStorage() {
	StaticModel = staticModel{}
	AgentModel = agentModel{}
	RouteTemplate = routeTemplate{}

	MuStatic = sync.RWMutex{}
	MuAgent = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
