package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(o *entities.Player, x float64, y float64) (*entities.Residence, error) {
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}

	r := Model.NewResidence(o, x, y)
	r.Name = "NoName"

	GenStepResidence(r)
	return r, nil
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(o *entities.Player, id uint) error {
	return TryRemove(o, entities.RESIDENCE, id, func(obj interface{}) {
		r := obj.(*entities.Residence)
		for _, h := range r.Targets {
			Model.Delete(h)
		}
		for _, s := range r.OutStep() {
			Model.Delete(s)
		}
		Model.Delete(r)
	})
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(o *entities.Player, x float64, y float64) (*entities.Company, error) {
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}

	c := Model.NewCompany(o, x, y)

	GenStepCompany(c)
	return c, nil
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(o *entities.Player, id uint) error {
	return TryRemove(o, entities.COMPANY, id, func(obj interface{}) {
		c := obj.(*entities.Company)
		for _, h := range c.Targets {
			Model.Delete(h)
		}
		for _, s := range c.InStep() {
			Model.Delete(s)
		}
		Model.Delete(c)
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
