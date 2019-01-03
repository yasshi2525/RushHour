package services

import (
	"fmt"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailNode create RailNode
func CreateRailNode(owner *entities.Player, x float64, y float64) *entities.RailNode {
	rn := entities.NewRailNode(GenID(entities.RAILNODE), owner, x, y)
	AddEntity(rn)
	return rn
}

// RemoveRailNode remove RailNode
func RemoveRailNode(owner *entities.Player, id uint) error {
	if rn, ok := Model.RailNodes[id]; ok {
		if ok, err := rn.CanRemove(); !ok {
			return err
		}
		if ok, err := IsAuth(owner, rn); !ok {
			return err
		}
		DelEntity(rn)
		return nil
	}
	revel.AppLog.Warnf("RailNode(%d) is already removed.", id)
	return nil
}

// IsAuth throws error when there is no permission
func IsAuth(owner *entities.Player, res entities.Ownable) (bool, error) {
	if res.Permits(owner) {
		return true, nil
	}
	return false, fmt.Errorf("no permission to operate %T: %+v", res, res)
}
