package services

import (
	"math/rand"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(x float64, y float64) (*entities.Residence, error) {
	r := entities.NewResidence(GenID(entities.RESIDENCE), x, y)
	r.Wait = Config.Residence.Interval.D.Seconds() * rand.Float64()
	r.Capacity = Config.Residence.Capacity
	r.Name = "NoName"
	AddEntity(r)
	GenStepResidence(r)
	return r, nil
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(id uint) error {
	return TryRemove(nil, entities.RESIDENCE, id, func(obj interface{}) {
		r := obj.(*entities.Residence)
		for _, h := range r.Targets {
			DelEntity(h)
		}
		for _, s := range r.Out() {
			DelEntity(s)
		}
		DelEntity(r)
	})
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(x float64, y float64) (*entities.Company, error) {
	c := entities.NewCompany(GenID(entities.COMPANY), x, y)
	c.Scale = Config.Company.Scale
	AddEntity(c)
	GenStepCompany(c)
	return c, nil
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(id uint) error {
	return TryRemove(nil, entities.COMPANY, id, func(obj interface{}) {
		c := obj.(*entities.Company)
		for _, h := range c.Targets {
			DelEntity(h)
		}
		for _, s := range c.In() {
			DelEntity(s)
		}
		DelEntity(c)
	})
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
