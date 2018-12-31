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

// StaticType represents type of resources for persistence
type StaticType string

// DynamicType represents type of resources not for persistence
type DynamicType uint

// StaticType represents type of resources indicating database table name
const (
	PLAYER    StaticType = "player"
	RESIDENCE            = "residence"
	COMPANY              = "company"
	RAILNODE             = "rail_node"
	RAILEDGE             = "rail_edge"
	STATION              = "station"
	GATE                 = "gate"
	PLATFORM             = "platform"
	LINE                 = "line"
	LINETASK             = "line_task"
	TRAIN                = "train"
	HUMAN                = "human"
)

// DynamicType represents type of resources not for persistence
const (
	STEP DynamicType = iota
)

// StaticTypes is list of all presistence resources.
// Dependecy order.
var StaticTypes []StaticType

// DynamicTypes is list of all non-persistence resources.
var DynamicTypes []DynamicType

// StaticInstances is list of each struc instance.
var StaticInstances []interface{}

// Dependecy order.
type staticModel struct {
	Players    map[uint]*entities.Player
	Residences map[uint]*entities.Residence
	Companies  map[uint]*entities.Company
	RailNodes  map[uint]*entities.RailNode
	RailEdges  map[uint]*entities.RailEdge
	Stations   map[uint]*entities.Station
	Gates      map[uint]*entities.Gate
	Platforms  map[uint]*entities.Platform
	Lines      map[uint]*entities.Line
	LineTasks  map[uint]*entities.LineTask
	Trains     map[uint]*entities.Train
	Humans     map[uint]*entities.Human
}

type dynamicModel struct {
	Steps map[uint]*entities.Step
}

type routeTemplate struct {
}

// Config defines game feature
var Config config

type nextID struct {
	Static  map[StaticType]*uint64
	Dynamic map[DynamicType]*uint64
}

// NextID has what number should be set
var NextID nextID

// Static is viewable feature including Step infomation.
var Static staticModel

// WillRemove represents the list of deleting in next Backup()
var WillRemove map[StaticType][]uint

// Dynamic is hidden feature and not be persisted.
var Dynamic dynamicModel

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
	StaticTypes = []StaticType{
		PLAYER,
		RESIDENCE,
		COMPANY,
		RAILNODE,
		RAILEDGE,
		STATION,
		GATE,
		PLATFORM,
		LINE,
		LINETASK,
		TRAIN,
		HUMAN,
	}
	DynamicTypes = []DynamicType{
		STEP,
	}

	StaticInstances = []interface{}{
		&entities.Player{},
		&entities.Residence{},
		&entities.Company{},
		&entities.RailNode{},
		&entities.RailEdge{},
		&entities.Station{},
		&entities.Gate{},
		&entities.Platform{},
		&entities.Line{},
		&entities.LineTask{},
		&entities.Train{},
		&entities.Human{},
	}

	Static = staticModel{
		Players:    make(map[uint]*entities.Player),
		Residences: make(map[uint]*entities.Residence),
		Companies:  make(map[uint]*entities.Company),
		RailNodes:  make(map[uint]*entities.RailNode),
		RailEdges:  make(map[uint]*entities.RailEdge),
		Stations:   make(map[uint]*entities.Station),
		Gates:      make(map[uint]*entities.Gate),
		Platforms:  make(map[uint]*entities.Platform),
		Lines:      make(map[uint]*entities.Line),
		LineTasks:  make(map[uint]*entities.LineTask),
		Trains:     make(map[uint]*entities.Train),
		Humans:     make(map[uint]*entities.Human),
	}

	Dynamic = dynamicModel{
		Steps: make(map[uint]*entities.Step),
	}
	RouteTemplate = routeTemplate{}

	WillRemove = make(map[StaticType][]uint)
	NextID = nextID{
		Static:  make(map[StaticType]*uint64),
		Dynamic: make(map[DynamicType]*uint64),
	}

	for _, t := range StaticTypes {
		WillRemove[t] = []uint{}
	}

	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}
