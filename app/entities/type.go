package entities

import "reflect"

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
	TRACK
	TRANSPORT
	STEP
)

// TypeList is list of ModelType
var TypeList []ModelType

// attribute represents meta information
type attribute struct {
	// Name represents identification of field
	Name string
	// Short represents short name description
	Short string
	// Table is table name corresponding to the field
	Table string
	// API is REST API name corresponding to the field
	API string
}

// attr represents attribute of each resource
var attr map[ModelType]*attribute
var types map[ModelType]reflect.Type
var nodes map[ModelType]bool
var edges map[ModelType]bool

// String returns identificable name
func (t ModelType) String() string {
	return attr[t].Name
}

func (t ModelType) Short() string {
	return attr[t].Short
}

// Table returns table name
func (t ModelType) Table() string {
	return attr[t].Table
}

// API returns REST API name
func (t ModelType) API() string {
	return attr[t].API
}

// Obj returns prototype pointer of instance
func (t ModelType) Obj(m *Model) Entity {
	obj := reflect.New(types[t])
	ptr := reflect.Indirect(obj).Addr().Interface()
	ptr.(Initializable).Init(m)
	return ptr.(Entity)
}

// Type returns type of field
func (t ModelType) Type() reflect.Type {
	return types[t]
}

// IsVisible returns whether user can refer it or not
func (t ModelType) IsVisible() bool {
	return attr[t].API != ""
}

// IsDB returns whether entity is persisted or not
func (t ModelType) IsDB() bool {
	return attr[t].Table != ""
}

func (t ModelType) IsRelayable() bool {
	return nodes[t]
}

func (t ModelType) IsConnectable() bool {
	return edges[t]
}

func InitType() {
	TypeList = []ModelType{
		PLAYER,
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
		TRACK,
		TRANSPORT,
		STEP,
	}

	attr = make(map[ModelType]*attribute)

	// name, short, table, api
	attr[PLAYER] = &attribute{"Player", "o", "players", "players"}
	attr[RESIDENCE] = &attribute{"Residence", "r", "residences", "residences"}
	attr[COMPANY] = &attribute{"Company", "c", "companies", "companies"}
	attr[RAILNODE] = &attribute{"RailNode", "rn", "rail_nodes", "rail_nodes"}
	attr[RAILEDGE] = &attribute{"RailEdge", "re", "rail_edges", "rail_edges"}
	attr[STATION] = &attribute{"Station", "st", "stations", "stations"}
	attr[GATE] = &attribute{"Gate", "g", "gates", "gates"}
	attr[PLATFORM] = &attribute{"Platform", "p", "platforms", "platforms"}
	attr[RAILLINE] = &attribute{"RailLine", "l", "rail_lines", "rail_lines"}
	attr[LINETASK] = &attribute{"LineTask", "lt", "line_tasks", "line_tasks"}
	attr[TRAIN] = &attribute{"Train", "t", "trains", "trains"}
	attr[HUMAN] = &attribute{"Human", "h", "humen", "humans"}

	attr[TRACK] = &attribute{"Track", "tk", "", ""}
	attr[TRANSPORT] = &attribute{"Transport", "x", "", ""}
	attr[STEP] = &attribute{"Step", "s", "", ""}

	types = make(map[ModelType]reflect.Type)
	nodes = make(map[ModelType]bool)
	edges = make(map[ModelType]bool)

	modelType := reflect.TypeOf(&Model{}).Elem()
	for idx, res := range TypeList {
		types[res] = modelType.Field(idx).Type.Elem().Elem()

		n := reflect.TypeOf((*Relayable)(nil)).Elem()
		e := reflect.TypeOf((*Connectable)(nil)).Elem()

		nodes[res] = types[res].Implements(n)
		edges[res] = types[res].Implements(e)
	}
}
