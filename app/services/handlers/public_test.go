package handlers

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/models"
)

func TestCreateResidence(t *testing.T) {
	models.Config.Residence.Capacity = 10
	models.Config.Residence.Interval = 1

	residence := CreateResidence(1, 1)

	if residence.Available > 1 {
		t.Errorf("Available should <= 1, but = %f", residence.Available)
	}
	if models.NextID.Residence != 1 {
		t.Errorf("NextID should 1, but = %d", models.NextID.Residence)
	}
}
