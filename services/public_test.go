package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
	"github.com/yasshi2525/RushHour/entities"
)

func TestCreateResidence(t *testing.T) {
	wd, _ := os.Getwd()
	conf, _ = config.Load(fmt.Sprintf("%s/../config", wd))
	conf.Game.Service.Routing.Worker = 1
	conf.Game.Entity.MaxScale = 2
	auther, _ = auth.GetAuther(conf.Secret.Auth)

	InitRepository()
	isInOperation = true

	admin, _ := CreatePlayer("test", "test", "test", 0, entities.Admin)

	residence, err := CreateResidence(admin, 1, 1)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveResidence(admin, residence.ID); err != nil {
		t.Error(err)
	}
}

func TestCreateCompany(t *testing.T) {
	wd, _ := os.Getwd()
	conf, _ = config.Load(fmt.Sprintf("%s/../config", wd))
	conf.Game.Service.Routing.Worker = 1
	conf.Game.Entity.MaxScale = 2
	auther, _ = auth.GetAuther(conf.Secret.Auth)

	InitRepository()
	isInOperation = true

	admin, _ := CreatePlayer("test", "test", "test", 0, entities.Admin)

	company, err := CreateCompany(admin, 1, 1)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveCompany(admin, company.ID); err != nil {
		t.Error(err)
	}
}

func TestCreateStep(t *testing.T) {
	wd, _ := os.Getwd()
	conf, _ = config.Load(fmt.Sprintf("%s/../config", wd))
	conf.Game.Service.Routing.Worker = 1
	conf.Game.Entity.MaxScale = 2
	auther, _ = auth.GetAuther(conf.Secret.Auth)

	InitRepository()
	isInOperation = true

	admin, _ := CreatePlayer("test", "test", "test", 0, entities.Admin)

	r, err := CreateResidence(admin, 1, 1)
	if err != nil {
		t.Error(err)
	}
	c, err := CreateCompany(admin, 2, 2)
	if err != nil {
		t.Error(err)
	}

	if got := len(r.OutSteps()); got != 1 {
		t.Errorf("Residence should be out 1, but %d", got)
	}
	if got := len(c.InSteps()); got != 1 {
		t.Errorf("Company should be in 1, but %d", got)
	}

	if err := RemoveResidence(admin, r.ID); err != nil {
		t.Errorf("RemoveResidence returns false, wanted true")
	}
	if err := RemoveCompany(admin, c.ID); err != nil {
		t.Errorf("RemoveCompany returns false, wanted true")
	}

	if got := len(Model.Steps); got != 0 {
		t.Errorf("Steps size should be 0, but %d", got)
	}
}
