package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// GenStepHuman Generate Step for Human.
// It's depend on where Human stay.
func GenStepHuman(h *entities.Human) {
	switch h.On {
	case entities.OnGround:
		// h - C for destination
		GenStep(h, h.To)
		// h -> G
		for _, g := range Model.Gates {
			GenStep(h, g)
		}
	case entities.OnPlatform:
		// h - G, P on Human
		GenStep(h, h.OnPlatform())
		GenStep(h, h.OnPlatform().WithGate)
	case entities.OnTrain:
		// do-nothing
	default:
		panic(fmt.Errorf("invalid type: %T %+v", h.On, h.On))
	}
}
