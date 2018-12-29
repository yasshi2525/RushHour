package handlers

import (
	"math/rand"
	"sync/atomic"

	"github.com/yasshi2525/RushHour/app/models"
	"github.com/yasshi2525/RushHour/app/models/entities"
)

// CreateResidence creates Residence and registers it to storage
func CreateResidence(x float64, y float64) entities.Residence {
	id := atomic.AddUint32(&models.NextID.Residence, 1)
	capacity := models.Config.Residence.Capacity
	available := models.Config.Residence.Interval * rand.Float64()

	residence := entities.Residence{
		Model:     entities.NewModel(id),
		Junction:  entities.NewJunction(x, y),
		Capacity:  capacity,
		Available: available,
		Targets:   []entities.Human{},
	}

	models.StaticModel.Residences = append(models.StaticModel.Residences, residence)

	return residence
}
