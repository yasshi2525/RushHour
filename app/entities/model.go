package entities

import (
	"reflect"
)

// ModelType represents type of resources
type ModelType uint

// ModelType represents type of resources
const (
	PLAYER ModelType = iota
	RESIDENCE
	COMPANY
	RAILNODE
	RAILEDGE
	STATION
	GATE
	PLATFORM
	RAILLINE
	LINETASK
	TRAIN
	HUMAN
	STEP
	AGENT
)

// MetaModel represents meta information of storage in memory
type MetaModel struct {
	// Attr represents attribute of each resource
	Attr map[ModelType]*Attribute
	// Map represents each resource map
	Map map[ModelType]reflect.Value
	// Type represents type of each resource
	Type map[ModelType]reflect.Type
	// List is list of ModelType
	List []ModelType
}

// Attribute represents meta information
type Attribute struct {
	// Name represents identification of field
	Name string
	// Short represents short name description
	Short string
	// Table is table name corresponding to the field
	Table string
	// API is REST API name corresponding to the field
	API string
}

// String returns identificable name
func (t ModelType) String() string {
	return Meta.Attr[t].Name
}

// Table returns table name
func (t ModelType) Table() string {
	return Meta.Attr[t].Table
}

// API returns REST API name
func (t ModelType) API() string {
	return Meta.Attr[t].API
}

// Obj returns prototype pointer of instance
func (t ModelType) Obj() interface{} {
	obj := reflect.New(Meta.Type[t].Elem())
	objptr := reflect.Indirect(obj).Addr().Interface()
	objptr.(Initializable).Init()
	return objptr
}

// Type returns type of field
func (t ModelType) Type() reflect.Type {
	return Meta.Type[t]
}

// IsVisible returns whether user can refer it or not
func (t ModelType) IsVisible() bool {
	return Meta.Attr[t].API != ""
}

// IsDB returns whether entity is persisted or not
func (t ModelType) IsDB() bool {
	return Meta.Attr[t].Table != ""
}

// Model represents data structure
type Model struct {
	Players    map[uint]*Player
	Residences map[uint]*Residence
	Companies  map[uint]*Company
	RailNodes  map[uint]*RailNode
	RailEdges  map[uint]*RailEdge
	Stations   map[uint]*Station
	Gates      map[uint]*Gate
	Platforms  map[uint]*Platform
	RailLines  map[uint]*RailLine
	LineTasks  map[uint]*LineTask
	Trains     map[uint]*Train
	Humans     map[uint]*Human

	Steps  map[uint]*Step
	Agents map[uint]*Agent

	NextIDs map[ModelType]*uint64
	// Remove represents the list of deleting in next Backup()
	Remove map[ModelType][]uint
}

// Meta is inner meta information for receivers.
var Meta *MetaModel

// InitModel initialize data structure
func InitModel() (*MetaModel, *Model) {
	meta := initMeta()
	Meta = meta
	model := initModel(meta)
	return meta, model
}

func initMeta() *MetaModel {
	meta := &MetaModel{
		make(map[ModelType]*Attribute),
		make(map[ModelType]reflect.Value),
		make(map[ModelType]reflect.Type),
		[]ModelType{PLAYER,
			RESIDENCE,
			COMPANY,
			RAILNODE,
			RAILEDGE,
			STATION,
			GATE,
			PLATFORM,
			RAILLINE,
			LINETASK,
			TRAIN,
			HUMAN,
			STEP,
			AGENT,
		},
	}

	// name, short, table, api
	meta.Attr[PLAYER] = &Attribute{"Player", "o", "players", "players"}
	meta.Attr[RESIDENCE] = &Attribute{"Residence", "r", "residences", "residences"}
	meta.Attr[COMPANY] = &Attribute{"Company", "c", "companies", "companies"}
	meta.Attr[RAILNODE] = &Attribute{"RailNode", "rn", "rail_nodes", "rail_nodes"}
	meta.Attr[RAILEDGE] = &Attribute{"RailEdge", "re", "rail_edges", "rail_edges"}
	meta.Attr[STATION] = &Attribute{"Station", "st", "stations", "stations"}
	meta.Attr[GATE] = &Attribute{"Gate", "g", "gates", "gates"}
	meta.Attr[PLATFORM] = &Attribute{"Platform", "p", "platforms", "platforms"}
	meta.Attr[RAILLINE] = &Attribute{"RailLine", "l", "rail_lines", "rail_lines"}
	meta.Attr[LINETASK] = &Attribute{"LineTask", "lt", "line_tasks", "line_tasks"}
	meta.Attr[TRAIN] = &Attribute{"Train", "t", "trains", "trains"}
	meta.Attr[HUMAN] = &Attribute{"Human", "h", "humen", "humans"}

	meta.Attr[STEP] = &Attribute{"Step", "s", "", ""}
	meta.Attr[AGENT] = &Attribute{"Agent", "a", "", ""}

	return meta
}

func initModel(meta *MetaModel) *Model {
	modelType := reflect.TypeOf(&Model{}).Elem()
	model := reflect.New(modelType).Elem()

	// set initialized map to field
	for idx, res := range meta.List {
		mapType := modelType.Field(idx).Type
		mapField := reflect.MakeMap(mapType)
		model.Field(idx).Set(mapField)

		// set type to metamodel
		meta.Map[res] = mapField
		meta.Type[res] = mapType.Elem()
	}

	obj := model.Addr().Interface().(*Model)
	obj.NextIDs = make(map[ModelType]*uint64)
	obj.Remove = make(map[ModelType][]uint)

	// set slice to specific fields
	for _, res := range meta.List {
		// NextID
		var id uint64
		obj.NextIDs[res] = &id
		// Remove
		if res.IsDB() {
			obj.Remove[res] = []uint{}
		}
	}

	return obj
}
