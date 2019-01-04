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
		GenWalkStep(h, h.To)
		// h -> G
		for _, g := range Model.Gates {
			GenWalkStep(h, g)
		}
	case entities.OnPlatform:
		// h - G, P on Human
		GenWalkStep(h, h.OnPlatform)
		GenWalkStep(h, h.OnPlatform.WithGate)
	case entities.OnTrain:
		// do-nothing
	default:
		panic(fmt.Errorf("invalid type: %T %+v", h.On, h.On))
	}
}
