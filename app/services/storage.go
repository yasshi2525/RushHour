package services

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
	validator "gopkg.in/go-playground/validator.v9"
)

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

// EntityType represents type of resources indicating database table name
type EntityType string

// EntityType represents type of resources indicating database table name
const (
	PLAYER    EntityType = "player"
	RESIDENCE            = "residence"
	COMPANY              = "company"
	RAILNODE             = "railnode"
	RAILEDGE             = "railedge"
	GATE                 = "gate"
	PLATFORM             = "platform"
	STATION              = "station"
	LINESTASK            = "linetask"
	LINE                 = "line"
	STEP                 = "step"
	TRAIN                = "train"
	HUMAN                = "human"
)

// EntityTypes is list of all entities.
var EntityTypes []EntityType

type staticModel struct {
	Players    map[uint]*entities.Player
	Residences map[uint]*entities.Residence
	Companies  map[uint]*entities.Company
	RailNodes  map[uint]*entities.RailNode
	RailEdges  map[uint]*entities.RailEdge
	Gates      map[uint]*entities.Gate
	Platforms  map[uint]*entities.Platform
	Stations   map[uint]*entities.Station
	LineTasks  map[uint]*entities.LineTask
	Lines      map[uint]*entities.Line
	Steps      map[uint]*entities.Step
	Trains     map[uint]*entities.Train
	Humen      map[uint]*entities.Human
}

type agentModel struct {
}

type routeTemplate struct {
}

// Config defines game feature
var Config config

// NextID has what number should be set
var NextID map[EntityType]*uint64

// Static is viewable feature including Step infomation.
var Static staticModel

// WillRemove represents the list of deleting in next Backup()
var WillRemove map[EntityType][]uint

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
	EntityTypes = []EntityType{
		PLAYER,
		COMPANY,
		RESIDENCE,
		RAILNODE,
		RAILEDGE,
		GATE,
		PLATFORM,
		STATION,
		STEP,
		TRAIN,
		HUMAN,
	}

	Static = staticModel{
		Players:    make(map[uint]*entities.Player),
		Companies:  make(map[uint]*entities.Company),
		Residences: make(map[uint]*entities.Residence),
		RailNodes:  make(map[uint]*entities.RailNode),
		RailEdges:  make(map[uint]*entities.RailEdge),
		Gates:      make(map[uint]*entities.Gate),
		Platforms:  make(map[uint]*entities.Platform),
		Stations:   make(map[uint]*entities.Station),
		LineTasks:  make(map[uint]*entities.LineTask),
		Lines:      make(map[uint]*entities.Line),
		Steps:      make(map[uint]*entities.Step),
		Trains:     make(map[uint]*entities.Train),
		Humen:      make(map[uint]*entities.Human),
	}
	Dynamic = agentModel{}
	RouteTemplate = routeTemplate{}

	WillRemove = make(map[EntityType][]uint)
	NextID = make(map[EntityType]*uint64)

	for _, t := range EntityTypes {
		WillRemove[t] = []uint{}
		var i uint64 = 1
		NextID[t] = &i
	}

	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
