package services

import (
	"fmt"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailNode create RailNode
func CreateRailNode(owner *entities.Player, x float64, y float64) *entities.RailNode {
	rn := entities.NewRailNode(GenID(entities.RAILNODE), owner, x, y)
	Repo.Static.RailNodes[rn.ID] = rn
	logOwnableNode(entities.RAILNODE, rn.ID, "created", rn.Loc, owner)
	return rn
}

// RemoveRailNode remove RailNode
func RemoveRailNode(owner *entities.Player, id uint) error {
	if rn, ok := Repo.Static.RailNodes[id]; ok {
		if in, out := len(rn.In), len(rn.Out); in > 0 || out > 0 {
			return fmt.Errorf("relations remain RailNode(%d)(in=%d, out=%d)", id, in, out)
		}
		if ok, err := IsAuth(owner, rn); !ok {
			return err
		}
		delete(Repo.Static.RailNodes, rn.ID)
		Repo.Static.WillRemove[entities.RAILNODE] = append(Repo.Static.WillRemove[entities.RAILNODE], id)
		logOwnableNode(entities.RAILNODE, id, "removed", rn.Loc, owner)
		return nil
	}
	revel.AppLog.Warnf("RailNode(%d) is already removed.", id)
	return nil
}

// IsAuth throws error when there is no permission
func IsAuth(owner *entities.Player, res entities.OwnableEntity) (bool, error) {
	if res.Permits(owner) {
		return true, nil
	}
	return false, fmt.Errorf("no permission to operate %T: %+v", res, res)
}

func logOwnableNode(res entities.StaticRes, id uint, op string, p *entities.Point, owner *entities.Player) {
	revel.AppLog.Infof("%s(%d) was %s at (%f, %f) by %s(%d)", res, id, op, p.X, p.Y, owner.LoginID, owner.ID)
}
