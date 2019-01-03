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

	residence, ok := CreateResidence(1, 1)
	if !ok {
		t.Errorf("CreateResidence returns false, wanted true")
	}

	if ok := RemoveResidence(residence.ID); !ok {
		t.Errorf("RemoveResidence returns false, wanted true")
	}
}

func TestCreateCompany(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	company, ok := CreateCompany(1, 1)
	if !ok {
		t.Errorf("CreateCompany returns false, wanted true")
	}

	if ok := RemoveCompany(company.ID); !ok {
		t.Errorf("RemoveCompany returns false, wanted true")
	}
}

func TestCreateStep(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	r, ok := CreateResidence(1, 1)
	if !ok {
		t.Errorf("CreateResidence returns false, wanted true")
	}
	c, ok := CreateCompany(2, 2)
	if !ok {
		t.Errorf("CreateCompany returns false, wanted true")
	}

	if got := len(r.Out()); got != 1 {
		t.Errorf("Residence should be out 1, but %d", got)
	}
	if got := len(c.In()); got != 1 {
		t.Errorf("Company should be in 1, but %d", got)
	}

	if ok := RemoveResidence(r.ID); !ok {
		t.Errorf("RemoveResidence returns false, wanted true")
	}
	if ok := RemoveCompany(c.ID); !ok {
		t.Errorf("RemoveCompany returns false, wanted true")
	}

	if got := len(Model.Steps); got != 0 {
		t.Errorf("Steps size should be 0, but %d", got)
	}
}
