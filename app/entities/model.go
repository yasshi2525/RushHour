package entities

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

const ZERO = 0

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
	Tracks     map[uint]*Track
	Transports map[uint]*Transport
	Steps      map[uint]*Step
	Agents     map[uint]*Agent

	NextIDs map[ModelType]*uint64
	// Deletes represents the list of deleting in next Backup()
	Deletes map[ModelType][]uint

	// Map represents each resource map
	Values map[ModelType]reflect.Value
}

func (m *Model) Find(res ModelType, idx uint) Indexable {
	if obj := m.Values[res].MapIndex(reflect.ValueOf(idx)); obj.IsValid() {
		return obj.Interface().(Indexable)
	}
	panic(fmt.Errorf("no corresponding object %v(%d)", res.Short(), idx))
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

func (m *Model) DeleteIf(o *Player, res ModelType, id uint) (Deletable, error) {
	raw := m.Values[res].MapIndex(reflect.ValueOf(id))
	// no id
	if !raw.IsValid() {
		return nil, fmt.Errorf("%v(%d) was already removed", res, id)
	}
	obj := raw.Interface().(Deletable)
	// no permission
	if !obj.Permits(o) {
		return obj, fmt.Errorf("no permission for %v to delete %v", o, obj)
	}
	// reference
	if err := obj.CheckDelete(); err != nil {
		return obj, err
	}
	obj.Delete()
	return obj, nil
}

func (m *Model) Delete(args ...UnReferable) {
	for _, obj := range args {
		obj.UnRef()
		m.Values[obj.Type()].SetMapIndex(
			reflect.ValueOf(obj.Idx()),
			reflect.Value{})
		m.Deletes[obj.Type()] = append(m.Deletes[obj.Type()], obj.Idx())
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
		if _, ok := res.Obj(m).(Relayable); ok {
			sum += m.Values[res].Len()
		}
	}
	return sum
}

func (m *Model) EdgeLen() int {
	var sum int
	for _, res := range TypeList {
		if _, ok := res.Obj(m).(Connectable); ok {
			sum += m.Values[res].Len()
		}
	}
	return sum
}

func (m *Model) DBLen() int {
	var sum int
	for _, res := range TypeList {
		if _, ok := res.Obj(m).(Persistable); ok {
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
	obj.Deletes = make(map[ModelType][]uint)
	obj.Values = make(map[ModelType]reflect.Value)

	// set slice to specific fields
	for idx, res := range TypeList {
		// NextID
		var id uint64
		obj.NextIDs[res] = &id
		// Deletes
		if res.IsDB() {
			obj.Deletes[res] = []uint{}
		}
		obj.Values[res] = model.Field(idx)
	}

	return obj
}
