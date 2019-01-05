package services

import (
	"fmt"
	"math/rand"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(o *entities.Player, x float64, y float64) (*entities.Residence, error) {
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}

	r := entities.NewResidence(GenID(entities.RESIDENCE), o, x, y)
	r.Wait = Config.Residence.Interval.D.Seconds() * rand.Float64()
	r.Capacity = Config.Residence.Capacity
	r.Name = "NoName"
	AddEntity(r)
	GenStepResidence(r)
	return r, nil
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(o *entities.Player, id uint) error {
	return TryRemove(o, entities.RESIDENCE, id, func(obj interface{}) {
		r := obj.(*entities.Residence)
		for _, h := range r.Targets {
			DelEntity(h)
		}
		for _, s := range r.OutStep() {
			DelEntity(s)
		}
		DelEntity(r)
	})
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(o *entities.Player, x float64, y float64) (*entities.Company, error) {
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}

	c := entities.NewCompany(GenID(entities.COMPANY), o, x, y)
	c.Scale = Config.Company.Scale
	AddEntity(c)
	GenStepCompany(c)
	return c, nil
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(o *entities.Player, id uint) error {
	return TryRemove(o, entities.COMPANY, id, func(obj interface{}) {
		c := obj.(*entities.Company)
		for _, h := range c.Targets {
			DelEntity(h)
		}
		for _, s := range c.InStep() {
			DelEntity(s)
		}
		DelEntity(c)
	})
}

// GenStepResidence generate Steps
func GenStepResidence(r *entities.Residence) {
	// R -> C
	for _, c := range Model.Companies {
		GenStep(r, c)
	}
	// R -> G
	for _, g := range Model.Gates {
		GenStep(r, g)
	}
}

// GenStepCompany generate Steps
func GenStepCompany(c *entities.Company) {
	// R -> C
	for _, r := range Model.Residences {
		GenStep(r, c)
	}
	// G -> C
	for _, g := range Model.Gates {
		GenStep(g, c)
	}
}
