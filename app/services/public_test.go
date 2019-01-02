package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yasshi2525/RushHour/app/entities"
)

func TestCreateResidence(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	var i uint64 = 1
	Static.NextIDs[entities.RESIDENCE] = &i

	residence := CreateResidence(1, 1)

	RemoveResidence(residence.ID)
}

func TestCreateCompany(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	var i uint64 = 1
	Static.NextIDs[entities.COMPANY] = &i

	RemoveCompany(CreateCompany(1, 1).ID)
}

func TestCreateStep(t *testing.T) {
	prev, _ := filepath.Abs(".")
	defer os.Chdir(prev)
	os.Chdir("../../")

	LoadConf()
	InitRepository()

	var i, j, k uint64 = 1, 1, 1
	Static.NextIDs[entities.RESIDENCE] = &i
	Static.NextIDs[entities.COMPANY] = &j
	Dynamic.NextIDs[entities.STEP] = &k

	r := CreateResidence(1, 1)
	c := CreateCompany(2, 2)

	if got := len(r.Out()); got != 1 {
		t.Errorf("Residence should be out 1, but %d", got)
	}
	if got := len(c.In()); got != 1 {
		t.Errorf("Company should be in 1, but %d", got)
	}

	RemoveResidence(r.ID)
	RemoveCompany(c.ID)

	if got := len(Dynamic.Steps); got != 0 {
		t.Errorf("Steps size should be 0, but %d", got)
	}
}
