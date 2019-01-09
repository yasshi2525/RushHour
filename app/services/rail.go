package services

import (
	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailNode create RailNode
func CreateRailNode(owner *entities.Player, x float64, y float64) (*entities.RailNode, error) {
	rn := Model.NewRailNode(owner, x, y)

	return rn, nil
}

// RemoveRailNode remove RailNode
func RemoveRailNode(owner *entities.Player, id uint) error {
	return TryRemove(owner, entities.RAILNODE, id, func(obj interface{}) {
		rn := obj.(*entities.RailNode)
		Model.Delete(rn)
	})
}

// ExtendRailNode extends Rail
func ExtendRailNode(o *entities.Player, from *entities.RailNode,
	x float64, y float64) (*entities.RailNode, *entities.RailEdge, *entities.RailEdge, error) {
	if err := CheckAuth(o, from); err != nil {
		return nil, nil, nil, err
	}
	to := Model.NewRailNode(from.Own, x, y)
	e1 := Model.NewRailEdge(from, to)
	e2 := Model.NewRailEdge(to, from)

	e1.Resolve(e2)
	e2.Resolve(e1)

	return to, e1, e2, nil
}
