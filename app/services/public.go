package services

import (
	"math/rand"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(x float64, y float64) *entities.Residence {
	r := entities.NewResidence(GenID(entities.RESIDENCE), x, y)
	r.Available = Config.Residence.Interval * rand.Float64()
	r.Capacity = Config.Residence.Capacity
	Repo.Static.Residences[r.ID] = r
	logNode(entities.RESIDENCE, r.ID, "created", r.Loc)

	GenStepResidence(r)
	return r
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(id uint) {
	if r, ok := Repo.Static.Residences[id]; ok {
		DelSteps(r.Out)
		delete(Repo.Static.Residences, id)
		Repo.Static.WillRemove[entities.RESIDENCE] = append(Repo.Static.WillRemove[entities.RESIDENCE], id)
		logNode(entities.RESIDENCE, id, "removed", r.Loc)
	} else {
		revel.AppLog.Warnf("%s(%d) is already removed.", entities.RESIDENCE, id)
		return
	}
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(x float64, y float64) *entities.Company {
	c := entities.NewCompany(GenID(entities.COMPANY), x, y)
	c.Scale = Config.Company.Scale
	Repo.Static.Companies[c.ID] = c
	logNode(entities.COMPANY, c.ID, "created", c.Loc)

	GenStepCompany(c)
	return c
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(id uint) {
	if c, ok := Repo.Static.Companies[id]; ok {
		DelSteps(c.In)
		delete(Repo.Static.Companies, id)
		Repo.Static.WillRemove[entities.COMPANY] = append(Repo.Static.WillRemove[entities.COMPANY], id)
		logNode(entities.COMPANY, id, "removed", c.Loc)
	} else {
		revel.AppLog.Warnf("%s(%d) is already removed.", entities.COMPANY, id)
		return
	}
}

func logNode(res entities.StaticRes, id uint, op string, p *entities.Point) {
	revel.AppLog.Infof("%s(%d) was %s at (%f, %f)", res, id, op, p.X, p.Y)
}

// GenStepResidence generate Steps
func GenStepResidence(r *entities.Residence) {
	// R -> C
	for _, c := range Repo.Static.Companies {
		GenStep(r, c, Config.Human.Weight)
	}
	// R -> G
	for _, g := range Repo.Static.Gates {
		GenStep(r, g, Config.Human.Weight)
	}
}

// GenStepCompany generate Steps
func GenStepCompany(c *entities.Company) {
	// R -> C
	for _, r := range Repo.Static.Residences {
		GenStep(r, c, Config.Human.Weight)
	}
	// G -> C
	for _, g := range Repo.Static.Gates {
		GenStep(g, c, Config.Human.Weight)
	}
}
