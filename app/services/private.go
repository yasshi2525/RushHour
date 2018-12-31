package services

import (
	"fmt"
	"sync/atomic"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailNode create RailNode
func CreateRailNode(owner *entities.Player, x float64, y float64) *entities.RailNode {
	id := uint(atomic.AddUint64(NextID.Static[RAILEDGE], 1))

	railNode := &entities.RailNode{
		Model:   entities.NewModel(id),
		Ownable: entities.NewOwnable(owner),
		Point:   entities.NewPoint(x, y),
		In:      []*entities.RailEdge{},
		Out:     []*entities.RailEdge{},
	}

	Static.RailNodes[id] = railNode
	logOwnableNode("RailNode", id, "created", &railNode.Point, owner)

	return railNode
}

// RemoveRailNode remove RailNode
func RemoveRailNode(owner *entities.Player, id uint) error {
	rn := Static.RailNodes[id]
	if rn == nil {
		revel.AppLog.Warnf("RailNode(%d) is already removed.", id)
		return nil
	}
	if owner.ID != rn.OwnerID {
		revel.AppLog.Warnf("%s try to remove RailNode(%d) owned by %s", owner.LoginID, rn.ID, rn.Owner.LoginID)
		return fmt.Errorf("no permission to remove RailNode(%d)", rn.ID)
	}
	delete(Static.RailNodes, rn.ID)
	WillRemove[RAILNODE] = append(WillRemove[RAILNODE], rn.ID)

	logOwnableNode("RailNode", id, "removed", &rn.Point, owner)
	return nil
}

func logOwnableNode(label string, id uint, op string, p *entities.Point, owner *entities.Player) {
	revel.AppLog.Infof("%s(%d) was %s at (%f, %f) by %s(%d)", label, id, op, p.X, p.Y, owner.LoginID, owner.ID)
}
