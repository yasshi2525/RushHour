package services

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/entities"
)

func TestCreateResidence(t *testing.T) {
	entities.InitStorage()
	entities.Config.Residence.Capacity = 10
	entities.Config.Residence.Interval = 1

	residence := CreateResidence(1, 1)

	if got := residence.Available; got > 1 {
		t.Errorf("Available should <= 1, but = %f", got)
	}

	RemoveResidence(residence)
}

func TestCreateCompany(t *testing.T) {
	entities.InitStorage()
	entities.Config.Company.Scale = 1

	RemoveCompany(CreateCompany(1, 1))
}

func TestCreateStep(t *testing.T) {
	entities.InitStorage()

	r := CreateResidence(1, 1)
	c := CreateCompany(2, 2)

	if got := len(r.Out); got != 1 {
		t.Errorf("Residence should be out 1, but %d", got)
	}
	if got := len(c.In); got != 1 {
		t.Errorf("Company should be in 1, but %d", got)
	}

	RemoveResidence(r)
	RemoveCompany(c)

	if got := len(entities.StaticModel.Steps); got != 0 {
		t.Errorf("Steps size should be 0, but %d", got)
	}
}
