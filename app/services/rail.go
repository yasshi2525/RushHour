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

// ExtendRailNode extends Rail
func ExtendRailNode(o *entities.Player, from *entities.RailNode,
	x float64, y float64) (*entities.RailNode, *entities.RailEdge, *entities.RailEdge, error) {
	if err := CheckAuth(o, from); err != nil {
		return nil, nil, nil, err
	}
	to := entities.NewRailNode(GenID(entities.RAILNODE), o, x, y)
	e1 := entities.NewRailEdge(GenID(entities.RAILEDGE), from, to)
	e2 := entities.NewRailEdge(GenID(entities.RAILEDGE), to, from)

	AddEntity(to, e1, e2)
	return to, e1, e2, nil
}
