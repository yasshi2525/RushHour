package services

import (
	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailNode create RailNode
func CreateRailNode(owner *entities.Player, x float64, y float64) (*entities.RailNode, error) {
	rn := entities.NewRailNode(GenID(entities.RAILNODE), owner, x, y)
	AddEntity(rn)
	return rn, nil
}

// RemoveRailNode remove RailNode
func RemoveRailNode(owner *entities.Player, id uint) error {
	return TryRemove(owner, entities.RAILNODE, id, func(obj interface{}) {
		rn := obj.(*entities.RailNode)
		DelEntity(rn)
	})
}
