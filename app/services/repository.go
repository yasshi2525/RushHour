package services

import (
	"sync"
	"sync/atomic"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// Static is viewable feature including Step infomation.
var Static *entities.StaticModel

// Dynamic is hidden feature and not be persisted.
var Dynamic *entities.DynamicModel

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
	Meta, Static, Dynamic = entities.InitGameMap()
	RouteTemplate = make(map[uint][]*entities.Node)
}

// GenID generate ID
func GenID(raw interface{}) uint {
	switch res := raw.(type) {
	case entities.StaticRes:
		return uint(atomic.AddUint64(Static.NextIDs[res], 1))
	case entities.DynamicRes:
		return uint(atomic.AddUint64(Dynamic.NextIDs[res], 1))
	default:
		revel.AppLog.Errorf("invalid type: %T: %+v", raw, raw)
		return 0
	}
}

// GenStep generate Step and resister it
func GenStep(from entities.Relayable, to entities.Relayable, weight float64) *entities.Step {
	s := entities.NewStep(GenID(entities.STEP), from, to, weight)
	Dynamic.Steps[s.ID] = s
	//logStep("created", s)
	return s
}

// DelSteps delete Step and unregister it
func DelSteps(steps map[uint]*entities.Step) {
	// be careful to change slice in range
	ids := []uint{}
	for _, s := range steps {
		ids = append(ids, s.ID)
	}
	for _, id := range ids {
		s := Dynamic.Steps[id]
		s.Unrelate()
		delete(Dynamic.Steps, s.ID)
		//logStep("removed", s)
	}
}

func logStep(op string, s *entities.Step) {
	from, to := s.From().Pos(), s.To().Pos()
	revel.AppLog.Debugf("Step(%d) was %s {%s => %s}", s.ID, op, from, to)
}
