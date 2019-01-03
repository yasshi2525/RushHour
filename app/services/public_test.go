package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateResidence(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	residence, err := CreateResidence(1, 1)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveResidence(residence.ID); err != nil {
		t.Error(err)
	}
}

func TestCreateCompany(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	company, err := CreateCompany(1, 1)
	if err != nil {
		t.Error(err)
	}

	if err := RemoveCompany(company.ID); err != nil {
		t.Error(err)
	}
}

func TestCreateStep(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	r, err := CreateResidence(1, 1)
	if err != nil {
		t.Error(err)
	}
	c, err := CreateCompany(2, 2)
	if err != nil {
		t.Error(err)
	}

	if got := len(r.Out()); got != 1 {
		t.Errorf("Residence should be out 1, but %d", got)
	}
	if got := len(c.In()); got != 1 {
		t.Errorf("Company should be in 1, but %d", got)
	}

	if err := RemoveResidence(r.ID); err != nil {
		t.Errorf("RemoveResidence returns false, wanted true")
	}
	if err := RemoveCompany(c.ID); err != nil {
		t.Errorf("RemoveCompany returns false, wanted true")
	}

	if got := len(Model.Steps); got != 0 {
		t.Errorf("Steps size should be 0, but %d", got)
	}
}
