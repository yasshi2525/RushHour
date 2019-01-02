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
	Static.Residences[r.ID] = r
	logNode(entities.RESIDENCE, r.ID, "created", r.Pos())

	GenStepResidence(r)
	return r
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(id uint) {
	if r, ok := Static.Residences[id]; ok {
		DelSteps(r.Out())
		delete(Static.Residences, id)
		Static.WillRemove[entities.RESIDENCE] = append(Static.WillRemove[entities.RESIDENCE], id)
		logNode(entities.RESIDENCE, id, "removed", r.Pos())
	} else {
		revel.AppLog.Warnf("%s(%d) is already removed.", entities.RESIDENCE, id)
		return
	}
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(x float64, y float64) *entities.Company {
	c := entities.NewCompany(GenID(entities.COMPANY), x, y)
	c.Scale = Config.Company.Scale
	Static.Companies[c.ID] = c
	logNode(entities.COMPANY, c.ID, "created", c.Pos())

	GenStepCompany(c)
	return c
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(id uint) {
	if c, ok := Static.Companies[id]; ok {
		DelSteps(c.In())
		delete(Static.Companies, id)
		Static.WillRemove[entities.COMPANY] = append(Static.WillRemove[entities.COMPANY], id)
		logNode(entities.COMPANY, id, "removed", c.Pos())
	} else {
		revel.AppLog.Warnf("%s(%d) is already removed.", entities.COMPANY, id)
		return
	}
}

func logNode(res entities.StaticRes, id uint, op string, p *entities.Point) {
	revel.AppLog.Infof("%s(%d) was %s at %s", res, id, op, p)
}

// GenStepResidence generate Steps
func GenStepResidence(r *entities.Residence) {
	// R -> C
	for _, c := range Static.Companies {
		GenStep(r, c, Config.Human.Weight)
	}
	// R -> G
	for _, g := range Static.Gates {
		GenStep(r, g, Config.Human.Weight)
	}
}

// GenStepCompany generate Steps
func GenStepCompany(c *entities.Company) {
	// R -> C
	for _, r := range Static.Residences {
		GenStep(r, c, Config.Human.Weight)
	}
	// G -> C
	for _, g := range Static.Gates {
		GenStep(g, c, Config.Human.Weight)
	}
}
