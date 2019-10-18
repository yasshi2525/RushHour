package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateResidence creates Residence and registers it to storage and step
func CreateResidence(o *entities.Player, x float64, y float64) (*entities.Residence, error) {
	if err := CheckMaintenance(); err != nil {
		return nil, err
	}
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}

	r := Model.NewResidence(x, y)
	r.Name = "NoName"

	AddOpLog("CreateResidence", o, r)
	return r, nil
}

// RemoveResidence remove Residence and related Step from storage
func RemoveResidence(o *entities.Player, id uint) error {
	if err := CheckMaintenance(); err != nil {
		return err
	}
	if r, err := Model.DeleteIf(o, entities.RESIDENCE, id); err != nil {
		return err
	} else {
		AddOpLog("RemoveResidence", o, r)
		return nil
	}
}

// CreateCompany creates Company and registers it to storage and step
func CreateCompany(o *entities.Player, x float64, y float64) (*entities.Company, error) {
	if err := CheckMaintenance(); err != nil {
		return nil, err
	}
	if o.Level != entities.Admin {
		return nil, fmt.Errorf("no permission")
	}
	c := Model.NewCompany(x, y)
	AddOpLog("CreateCompany", o, c)
	return c, nil
}

// RemoveCompany remove Company and related Step from storage
func RemoveCompany(o *entities.Player, id uint) error {
	if err := CheckMaintenance(); err != nil {
		return err
	}
	if c, err := Model.DeleteIf(o, entities.COMPANY, id); err != nil {
		return err
	} else {
		AddOpLog("RemoveCompany", o, c)
		return nil
	}
}
