package entities

import (
	"reflect"
	"sync/atomic"
)

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
	Steps      map[uint]*Step
	Agents     map[uint]*Agent

	NextIDs map[ModelType]*uint64
	// Remove represents the list of deleting in next Backup()
	Remove map[ModelType][]uint

	// Map represents each resource map
	Values map[ModelType]reflect.Value
}

func (m *Model) Find(res ModelType, idx uint) Indexable {
	return m.Values[res].MapIndex(reflect.ValueOf(idx)).Interface().(Indexable)
}

func (m *Model) ForEach(res ModelType, callback func(Indexable)) {
	mapdata := m.Values[res]
	for _, key := range mapdata.MapKeys() {
		callback(mapdata.MapIndex(key).Interface().(Indexable))
	}
}

func (m *Model) GenID(res ModelType) uint {
	return uint(atomic.AddUint64(m.NextIDs[res], 1))
}

func (m *Model) Add(args ...Indexable) {
	for _, obj := range args {
		m.Values[obj.Type()].SetMapIndex(
			reflect.ValueOf(obj.Idx()),
			reflect.ValueOf(obj))
	}
}

func (m *Model) Delete(args ...UnReferable) {
	for _, obj := range args {
		obj.UnRef()
		m.Values[obj.Type()].SetMapIndex(
			reflect.ValueOf(obj.Idx()),
			reflect.Value{})
		m.Remove[obj.Type()] = append(m.Remove[obj.Type()], obj.Idx())
	}
}

func (m *Model) Ids(res ModelType) []uint {
	ids := make([]uint, m.Values[res].Len())
	var i int
	for _, key := range m.Values[res].MapKeys() {
		ids[i] = uint(key.Uint())
		i++
	}
	return ids
}

func (m *Model) Len() int {
	var sum int
	for _, res := range TypeList {
		sum += m.Values[res].Len()
	}
	return sum
}

func (m *Model) NodeLen() int {
	var sum int
	for _, res := range TypeList {
		if _, ok := res.Obj().(Relayable); ok {
			sum += m.Values[res].Len()
		}
	}
	return sum
}

func (m *Model) EdgeLen() int {
	var sum int
	for _, res := range TypeList {
		if _, ok := res.Obj().(Connectable); ok {
			sum += m.Values[res].Len()
		}
	}
	return sum
}

func (m *Model) DBLen() int {
	var sum int
	for _, res := range TypeList {
		if _, ok := res.Obj().(Persistable); ok {
			sum += m.Values[res].Len()
		}
	}
	return sum
}

func NewModel() *Model {
	modelType := reflect.TypeOf(&Model{}).Elem()
	model := reflect.New(modelType).Elem()

	// set initialized map to field
	for idx := range TypeList {
		mapType := modelType.Field(idx).Type
		mapField := reflect.MakeMap(mapType)
		model.Field(idx).Set(mapField)
	}

	obj := model.Addr().Interface().(*Model)
	obj.NextIDs = make(map[ModelType]*uint64)
	obj.Remove = make(map[ModelType][]uint)
	obj.Values = make(map[ModelType]reflect.Value)

	// set slice to specific fields
	for idx, res := range TypeList {
		// NextID
		var id uint64
		obj.NextIDs[res] = &id
		// Remove
		if res.IsDB() {
			obj.Remove[res] = []uint{}
		}
		obj.Values[res] = model.Field(idx)
	}

	return obj
}
