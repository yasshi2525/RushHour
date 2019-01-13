package services

import (
	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/route"
)

// CreateRailNode create RailNode
func CreateRailNode(o *entities.Player, x float64, y float64) (*entities.RailNode, error) {
	rn := Model.NewRailNode(o, x, y)
	AddOpLog("CreateRailNode", o, rn)
	return rn, nil
}

// RemoveRailNode remove RailNode
func RemoveRailNode(o *entities.Player, id uint) error {
	if rn, err := Model.DeleteIf(o, entities.RAILNODE, id); err != nil {
		return err
	} else {
		AddOpLog("RemoveRailNode", o, rn)
		return nil
	}
}

// ExtendRailNode extends Rail
func ExtendRailNode(o *entities.Player, from *entities.RailNode,
	x float64, y float64) (*entities.RailNode, *entities.RailEdge, *entities.RailEdge, error) {
	if err := CheckAuth(o, from); err != nil {
		return nil, nil, nil, err
	}
	to, e1 := from.Extend(x, y)
	route.RefreshTracks(o, Const.Routing.Worker)
	for _, l := range from.RailLines {
		if l.ReRouting {
			route.RefreshTransports(l, Const.Routing.Worker)
		}
	}
	AddOpLog("ExtendRailNode", o, from, to, e1, e1.Reverse)
	return to, e1, e1.Reverse, nil
}

// RemoveRailEdge remove RailEdge
func RemoveRailEdge(o *entities.Player, id uint) error {
	if re, err := Model.DeleteIf(o, entities.RAILEDGE, id); err != nil {
		return err
	} else {
		AddOpLog("RemoveRailEdge", o, re)
		return nil
	}
}
