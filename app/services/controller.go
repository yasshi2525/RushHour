package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// ViewDelegateMap returns delegate Entity for client view.
func ViewDelegateMap(x float64, y float64, scale float64, delegate float64) *entities.DelegateMap {
	dm := &entities.DelegateMap{}
	dm.Init()
	Model.RootCluster.ViewMap(dm, x, y, scale, delegate)
	return dm
}

// CheckAuth throws error when there is no permission
func CheckAuth(owner *entities.Player, res entities.Entity) error {
	if res.B().Permits(owner) {
		return nil
	}
	return fmt.Errorf("no permission to operate %v", res)
}

func CheckArea(x float64, y float64) error {
	if (&entities.Point{X: x, Y: y}).IsIn(0, 0, serviceConf.AppConf.Game.Entity.MaxScale) {
		return nil
	}
	return fmt.Errorf("out of bounds")
}
