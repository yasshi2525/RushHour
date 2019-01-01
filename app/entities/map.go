package entities

import (
	"reflect"
)

// StaticRes represents type of resources for persistence
type StaticRes uint

// DynamicRes represents type of resources not for persistence
type DynamicRes uint

// StaticRes represents type of resources for persistence
// DynamicRes represents type of resources not for persistence
const (
	PLAYER StaticRes = iota
	RESIDENCE
	COMPANY
	RAILNODE
	RAILEDGE
	STATION
	GATE
	PLATFORM
	LINE
	LINETASK
	TRAIN
	HUMAN
	STEP DynamicRes = iota
)

// MetaModel represents meta information of storage in memory
type MetaModel struct {
	// Static is data storage
	Static map[StaticRes]*MetaStatic
	// StaticType represents type of field
	StaticType map[StaticRes]reflect.Type
	// StaticValue represents value of field
	StaticValue map[StaticRes]reflect.Value
	// StaticList is list of StaticRes
	StaticList []StaticRes

	// Dynamic is data storage
	Dynamic map[DynamicRes]*MetaDynamic
	// DynamicType  represents type of fi
	DynamicType map[DynamicRes]reflect.Type
	// DynamicValue represents value of field
	DynamicValue map[DynamicRes]reflect.Value
	// DynamicList is data storage
	DynamicList []DynamicRes
}

// MetaStatic represents meta information of Static
type MetaStatic struct {
	// Name represents identification of field
	Name string
	// Short represents short name description
	Short string
	// Table is table name corresponding to the field of Static
	Table string
	// API is REST API name corresponding to the field of Static
	API string
}

// MetaDynamic represents meta information of Dynamic
type MetaDynamic struct {
	// Name represents identification of field
	Name string
	// Short represents short name description
	Short string
}

// String returns identificable name
func (t StaticRes) String() string {
	return Meta.Static[t].Name
}

// Table returns table name
func (t StaticRes) Table() string {
	return Meta.Static[t].Table
}

// API returns REST API name
func (t StaticRes) API() string {
	return Meta.Static[t].API
}

// Obj returns prototype pointer of instance
func (t StaticRes) Obj() interface{} {
	return reflect.New(Meta.StaticType[t].Elem().Elem()).Elem().Addr().Interface()
}

// Type returns type of field
func (t StaticRes) Type() reflect.Type {
	return Meta.StaticType[t]
}

// StaticModel represents data structure for view
type StaticModel struct {
	Players    map[uint]*Player    `json:"players"`
	Residences map[uint]*Residence `json:"residences"`
	Companies  map[uint]*Company   `json:"companies"`
	RailNodes  map[uint]*RailNode  `json:"rail_nodes"`
	RailEdges  map[uint]*RailEdge  `json:"rail_edges"`
	Stations   map[uint]*Station   `json:"stations"`
	Gates      map[uint]*Gate      `json:"gates"`
	Platforms  map[uint]*Platform  `json:"platforms"`
	RailLines  map[uint]*RailLine  `json:"rail_lines"`
	LineTasks  map[uint]*LineTask  `json:"line_tasks"`
	Trains     map[uint]*Train     `json:"trains"`
	Humans     map[uint]*Human     `json:"humans"`

	NextIDs map[StaticRes]*uint64
	// WillRemove represents the list of deleting in next Backup()
	WillRemove map[StaticRes][]uint
}

// DynamicModel represents data structure for agent
type DynamicModel struct {
	Steps map[uint]*Step

	NextIDs map[DynamicRes]*uint64
}

// Meta is inner meta information for receivers.
var Meta *MetaModel

// InitGameMap initialize data structure
func InitGameMap() (*MetaModel, *StaticModel, *DynamicModel) {
	meta := initMetaModel()
	Meta = meta
	static, dynamic := initModel(meta)
	resolveModel(meta, static, dynamic)
	return meta, static, dynamic
}

func initMetaModel() *MetaModel {
	meta := &MetaModel{
		make(map[StaticRes]*MetaStatic),
		make(map[StaticRes]reflect.Type),
		make(map[StaticRes]reflect.Value),
		[]StaticRes{PLAYER,
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
			HUMAN},
		make(map[DynamicRes]*MetaDynamic),
		make(map[DynamicRes]reflect.Type),
		make(map[DynamicRes]reflect.Value),
		[]DynamicRes{STEP},
	}

	// name, short, table, api, instance, next_id
	meta.Static[PLAYER] = &MetaStatic{"Player", "o", "players", "players"}
	meta.Static[RESIDENCE] = &MetaStatic{"Residence", "r", "residences", "residences"}
	meta.Static[COMPANY] = &MetaStatic{"Company", "c", "companies", "companies"}
	meta.Static[RAILNODE] = &MetaStatic{"RailNode", "rn", "rail_nodes", "rail_nodes"}
	meta.Static[RAILEDGE] = &MetaStatic{"RailEdge", "re", "rail_edges", "rail_edges"}
	meta.Static[STATION] = &MetaStatic{"Station", "st", "stations", "stations"}
	meta.Static[GATE] = &MetaStatic{"Gate", "g", "gates", "gates"}
	meta.Static[PLATFORM] = &MetaStatic{"Platform", "p", "platforms", "platforms"}
	meta.Static[LINE] = &MetaStatic{"RailLine", "l", "rail_lines", "rail_lines"}
	meta.Static[LINETASK] = &MetaStatic{"LineTask", "lt", "line_tasks", "line_tasks"}
	meta.Static[TRAIN] = &MetaStatic{"Train", "t", "trains", "trains"}
	meta.Static[HUMAN] = &MetaStatic{"Human", "h", "humen", "humans"}

	return meta
}

func initModel(meta *MetaModel) (*StaticModel, *DynamicModel) {
	static := &StaticModel{
		make(map[uint]*Player),
		make(map[uint]*Residence),
		make(map[uint]*Company),
		make(map[uint]*RailNode),
		make(map[uint]*RailEdge),
		make(map[uint]*Station),
		make(map[uint]*Gate),
		make(map[uint]*Platform),
		make(map[uint]*RailLine),
		make(map[uint]*LineTask),
		make(map[uint]*Train),
		make(map[uint]*Human),
		make(map[StaticRes]*uint64),
		make(map[StaticRes][]uint),
	}

	dynamic := &DynamicModel{
		make(map[uint]*Step),
		make(map[DynamicRes]*uint64),
	}

	for _, res := range meta.StaticList {
		static.WillRemove[res] = []uint{}
		var id uint64 = 1
		static.NextIDs[res] = &id
	}

	for _, res := range meta.DynamicList {
		var id uint64 = 1
		dynamic.NextIDs[res] = &id
	}

	return static, dynamic
}

func resolveModel(meta *MetaModel, static *StaticModel, dynamic *DynamicModel) {
	meta.StaticType[PLAYER] = reflect.TypeOf(static.Players)
	meta.StaticType[RESIDENCE] = reflect.TypeOf(static.Residences)
	meta.StaticType[COMPANY] = reflect.TypeOf(static.Companies)
	meta.StaticType[RAILNODE] = reflect.TypeOf(static.RailNodes)
	meta.StaticType[RAILEDGE] = reflect.TypeOf(static.RailEdges)
	meta.StaticType[STATION] = reflect.TypeOf(static.Stations)
	meta.StaticType[GATE] = reflect.TypeOf(static.Gates)
	meta.StaticType[PLATFORM] = reflect.TypeOf(static.Platforms)
	meta.StaticType[LINE] = reflect.TypeOf(static.RailLines)
	meta.StaticType[LINETASK] = reflect.TypeOf(static.LineTasks)
	meta.StaticType[TRAIN] = reflect.TypeOf(static.Trains)
	meta.StaticType[HUMAN] = reflect.TypeOf(static.Humans)

	meta.StaticValue[PLAYER] = reflect.ValueOf(static.Players)
	meta.StaticValue[RESIDENCE] = reflect.ValueOf(static.Residences)
	meta.StaticValue[COMPANY] = reflect.ValueOf(static.Companies)
	meta.StaticValue[RAILNODE] = reflect.ValueOf(static.RailNodes)
	meta.StaticValue[RAILEDGE] = reflect.ValueOf(static.RailEdges)
	meta.StaticValue[STATION] = reflect.ValueOf(static.Stations)
	meta.StaticValue[GATE] = reflect.ValueOf(static.Gates)
	meta.StaticValue[PLATFORM] = reflect.ValueOf(static.Platforms)
	meta.StaticValue[LINE] = reflect.ValueOf(static.RailLines)
	meta.StaticValue[LINETASK] = reflect.ValueOf(static.LineTasks)
	meta.StaticValue[TRAIN] = reflect.ValueOf(static.Trains)
	meta.StaticValue[HUMAN] = reflect.ValueOf(static.Humans)

	meta.DynamicType[STEP] = reflect.TypeOf(dynamic.Steps)
	meta.DynamicValue[STEP] = reflect.ValueOf(dynamic.Steps)
}
