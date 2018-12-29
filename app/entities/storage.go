package entities

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/revel/revel"
	validator "gopkg.in/go-playground/validator.v9"
)

type nextID struct {
	Residence uint64
	Company   uint64
	Step      uint64
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
	Residences map[uint]*Residence
	Companies  map[uint]*Company
	Gates      map[uint]*Gate
	Platforms  map[uint]*Platform
	Train      map[uint]*Train
	Steps      map[uint]*Step
}

type agentModel struct {
}

type routeTemplate struct {
}

// Config defines game feature
var Config config

// NextID has what number should be set
var NextID nextID

// StaticModel is viewable feature including Step infomation.
var StaticModel staticModel

// AgentModel is hidden feature and not be persisted/
var AgentModel agentModel

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate routeTemplate

// MuStatic is mutex lock for StaticModel
var MuStatic sync.RWMutex

// MuAgent is mutex lock for AgentModel
var MuAgent sync.RWMutex

// MuRoute is mutex lock for routing
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
	StaticModel = staticModel{
		Companies:  make(map[uint]*Company),
		Residences: make(map[uint]*Residence),
		Gates:      make(map[uint]*Gate),
		Platforms:  make(map[uint]*Platform),
		Steps:      make(map[uint]*Step),
	}
	AgentModel = agentModel{}
	RouteTemplate = routeTemplate{}

	MuStatic = sync.RWMutex{}
	MuAgent = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
