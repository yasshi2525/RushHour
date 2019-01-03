package services

import (
	"math/rand"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(x float64, y float64) (*entities.Residence, bool) {
	r := entities.NewResidence(GenID(entities.RESIDENCE), x, y)
	r.Wait = Config.Residence.Interval.Duration.Seconds() * rand.Float64()
	r.Capacity = Config.Residence.Capacity
	r.Name = "NoName"
	AddEntity(r)
	GenStepResidence(r)
	return r, true
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(id uint) bool {
	if r, ok := Model.Residences[id]; ok {
		for _, h := range r.Targets {
			DelEntity(h)
		}
		for _, s := range r.Out() {
			DelEntity(s)
		}
		DelEntity(r)
		return true
	}
	revel.AppLog.Warnf("%s(%d) is already removed.", entities.RESIDENCE, id)
	return false

}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(x float64, y float64) (*entities.Company, bool) {
	c := entities.NewCompany(GenID(entities.COMPANY), x, y)
	c.Scale = Config.Company.Scale
	AddEntity(c)
	GenStepCompany(c)
	return c, true
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(id uint) bool {
	if c, ok := Model.Companies[id]; ok {
		for _, h := range c.Targets {
			DelEntity(h)
		}
		for _, s := range c.In() {
			DelEntity(s)
		}
		DelEntity(c)
		return true
	}
	revel.AppLog.Warnf("%s(%d) is already removed.", entities.COMPANY, id)
	return false
}

// GenStepResidence generate Steps
func GenStepResidence(r *entities.Residence) {
	// R -> C
	for _, c := range Model.Companies {
		GenStep(r, c, Config.Human.Weight)
	}
	// R -> G
	for _, g := range Model.Gates {
		GenStep(r, g, Config.Human.Weight)
	}
}

// GenStepCompany generate Steps
func GenStepCompany(c *entities.Company) {
	// R -> C
	for _, r := range Model.Residences {
		GenStep(r, c, Config.Human.Weight)
	}
	// G -> C
	for _, g := range Model.Gates {
		GenStep(g, c, Config.Human.Weight)
	}
}
