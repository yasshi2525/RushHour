package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// Generate Step for Human.
// It's depend on where Human stay.
func GenStepHuman(h *entities.Human) {
	w := Config.Human.Weight
	switch h.On {
	case entities.OnGround:
		// h - C for destination
		GenStep(h, h.To, w)
		// h -> G
		for _, g := range Repo.Static.Gates {
			GenStep(h, g, w)
		}
	case entities.OnPlatform:
		// h - G, P on Human
		GenStep(h, h.OnPlatform, w)
		GenStep(h, h.OnPlatform.InStation.Gate, w)
	case entities.OnTrain:
		// do-nothing
	default:
		panic(fmt.Errorf("invalid type: %T %+v", h.On, h.On))
	}
}
