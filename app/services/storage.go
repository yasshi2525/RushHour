package services

import (
	"sync"
	"sync/atomic"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

type routeTemplate struct {
}

// Repository has all data of game
type Repository struct {
	// Static is viewable feature including Step infomation.
	Static *entities.StaticModel
	// Dynamic is hidden feature and not be persisted.
	Dynamic *entities.DynamicModel
	// Meta represents meta information of data structure
	Meta *entities.MetaModel
}

// Repo has all data of game
var Repo *Repository

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate routeTemplate

// MuStatic is mutex lock for Static
var MuStatic sync.RWMutex

// MuDynamic is mutex lock for Dynamic
var MuDynamic sync.RWMutex

// MuRoute is mutex lock for routing
var MuRoute sync.Mutex

// InitStorage initialize storage
func InitStorage() {
	m, s, d := entities.InitGameMap()
	Repo = &Repository{
		Meta:    m,
		Static:  s,
		Dynamic: d,
	}

	RouteTemplate = routeTemplate{}

	MuStatic = sync.RWMutex{}
	MuDynamic = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}

// GenID generate ID
func GenID(raw interface{}) uint {
	switch res := raw.(type) {
	case entities.StaticRes:
		return uint(atomic.AddUint64(Repo.Static.NextIDs[res], 1))
	case entities.DynamicRes:
		return uint(atomic.AddUint64(Repo.Dynamic.NextIDs[res], 1))
	default:
		revel.AppLog.Errorf("invalid type: %T: %+v", raw, raw)
		return 0
	}
}

// GenStep generate Step and resister it
func GenStep(from entities.Junction, to entities.Junction, weight float64) *entities.Step {
	s := entities.NewStep(GenID(entities.STEP), from, to, weight)
	Repo.Dynamic.Steps[s.ID] = s
	logStep("created", s)
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
		s := Repo.Dynamic.Steps[id]
		s.Unrelate()
		delete(Repo.Dynamic.Steps, s.ID)
		logStep("removed", s)
	}
}

func logStep(op string, s *entities.Step) {
	revel.AppLog.Debugf("Step(%d) was %s {(%f, %f) => (%f, %f)}",
		s.ID, op, s.From.Pos().X, s.From.Pos().Y, s.To.Pos().X, s.To.Pos().Y)
}
