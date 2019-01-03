package services

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// Model is contained all data for gaming
var Model *entities.Model

// Meta represents meta information of data structure
var Meta *entities.MetaModel

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate map[uint][]*entities.Node

// MuStatic is mutex lock for Static
var MuStatic sync.RWMutex

// MuDynamic is mutex lock for Dynamic
var MuDynamic sync.RWMutex

// MuRoute is mutex lock for routing
var MuRoute sync.Mutex

// InitLock must prepare first.
func InitLock() {
	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}

// InitRepository initialize storage
func InitRepository() {
	Meta, Model = entities.InitModel()
	RouteTemplate = make(map[uint][]*entities.Node)
}

// GenID generate ID
func GenID(t entities.ModelType) uint {
	return uint(atomic.AddUint64(Model.NextIDs[t], 1))
}

// GenStep generate Step and resister it
func GenStep(from entities.Relayable, to entities.Relayable, weight float64) *entities.Step {
	s := entities.NewStep(GenID(entities.STEP), from, to, weight)
	AddEntity(s)
	return s
}

// AddEntity create entity and register to map
func AddEntity(obj entities.Indexable) {
	Meta.Map[obj.Type()].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.ValueOf(obj))
	revel.AppLog.Debugf("created %v", obj)
}

// DelEntity unrefer obj and delete it from map
func DelEntity(obj entities.UnReferable) {
	obj.UnRef()
	Meta.Map[obj.Type()].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.Value{})
	Model.Remove[obj.Type()] = append(Model.Remove[obj.Type()], obj.Idx())
	revel.AppLog.Debugf("removed %v", obj)
}

// ForeachModel iterates specified map
func ForeachModel(res entities.ModelType, callback func(interface{})) {
	mapdata := Meta.Map[res]
	for _, key := range mapdata.MapKeys() {
		callback(mapdata.MapIndex(key).Interface())
	}
}
