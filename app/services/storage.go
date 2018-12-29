package services

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
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
	Residences map[uint]*entities.Residence
	Companies  map[uint]*entities.Company
	Gates      map[uint]*entities.Gate
	Platforms  map[uint]*entities.Platform
	Train      map[uint]*entities.Train
	Steps      map[uint]*entities.Step
}

type agentModel struct {
}

type routeTemplate struct {
}

// Config defines game feature
var Config config

// NextID has what number should be set
var NextID nextID

// Static is viewable feature including Step infomation.
var Static staticModel

// Dynamic is hidden feature and not be persisted/
var Dynamic agentModel

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate routeTemplate

// MuStatic is mutex lock for Static
var MuStatic sync.RWMutex

// MuDynamic is mutex lock for Dynamic
var MuDynamic sync.RWMutex

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
	Static = staticModel{
		Companies:  make(map[uint]*entities.Company),
		Residences: make(map[uint]*entities.Residence),
		Gates:      make(map[uint]*entities.Gate),
		Platforms:  make(map[uint]*entities.Platform),
		Steps:      make(map[uint]*entities.Step),
	}
	Dynamic = agentModel{}
	RouteTemplate = routeTemplate{}

	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
