package services

import (
	"math/rand"
	"sync/atomic"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(x float64, y float64) *entities.Residence {
	id := uint(atomic.AddUint64(NextID[RESIDENCE], 1))
	capacity := Config.Residence.Capacity
	available := Config.Residence.Interval * rand.Float64()

	residence := &entities.Residence{
		Model:     entities.NewModel(id),
		Junction:  entities.NewJunction(x, y),
		Capacity:  capacity,
		Available: available,
		Targets:   []entities.Human{},
	}

	Static.Residences[id] = residence
	logNode("Residence", id, "created", &residence.Point)

	for _, c := range Static.Companies {
		createStep(&residence.Junction, &c.Junction, 1.0)
	}
	for _, g := range Static.Gates {
		createStep(&residence.Junction, &g.Junction, 1.0)
	}

	return residence
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(id uint) {
	r := Static.Residences[id]
	if r == nil {
		revel.AppLog.Warnf("residence(%d) is already removed.", id)
		return
	}

	for _, s := range r.Out {
		delete(Static.Steps, s.ID)
		logStep("removed", s)
	}
	delete(Static.Residences, r.ID)
	WillRemove[RESIDENCE] = append(WillRemove[RESIDENCE], r.ID)
	logNode("Residence", r.ID, "removed", &r.Point)
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(x float64, y float64) *entities.Company {
	id := uint(atomic.AddUint64(NextID[COMPANY], 1))
	scale := Config.Company.Scale

	company := &entities.Company{
		Model:    entities.NewModel(id),
		Junction: entities.NewJunction(x, y),
		Scale:    scale,
		Targets:  []entities.Human{},
	}

	Static.Companies[id] = company
	logNode("Company", id, "created", &company.Point)

	for _, r := range Static.Residences {
		createStep(&r.Junction, &company.Junction, 1.0)
	}
	for _, g := range Static.Gates {
		createStep(&g.Junction, &company.Junction, 1.0)
	}

	return company
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(id uint) {
	c := Static.Companies[id]
	if c == nil {
		revel.AppLog.Warnf("company(%d) is already removed.", id)
		return
	}
	for _, s := range c.In {
		delete(Static.Steps, s.ID)
		logStep("removed", s)
	}
	delete(Static.Companies, c.ID)
	WillRemove[COMPANY] = append(WillRemove[COMPANY], c.ID)
	logNode("Company", c.ID, "removed", &c.Point)
}

func createStep(from *entities.Junction, to *entities.Junction, weight float64) *entities.Step {
	id := uint(atomic.AddUint64(NextID[STEP], 1))
	step := &entities.Step{
		ID:     id,
		From:   from,
		To:     to,
		Weight: weight,
	}
	from.Out = append(from.Out, step)
	to.In = append(to.In, step)
	Static.Steps[id] = step
	logStep("created", step)
	return step
}

func logNode(label string, id uint, op string, p *entities.Point) {
	revel.AppLog.Infof("%s(%d) was %s at (%f, %f)", label, id, op, p.X, p.Y)
}

func logStep(op string, s *entities.Step) {
	revel.AppLog.Debugf("Step(%d) was %s {(%f, %f) => (%f, %f)}",
		s.ID, op, s.From.X, s.From.Y, s.To.X, s.To.Y)
}
