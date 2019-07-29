package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

// RefreshTracks set minimum distance trace on specified rail.
func RefreshTracks(o *entities.Player, parallel int) map[uint]*Model {
	o.ClearTracks()
	template := scanRail(o)

	payload, _ := Search(context.Background(), entities.RAILNODE, parallel, template)

	for destID, model := range payload.Route {
		for deptID, dept := range model.Nodes[entities.RAILNODE] {
			if dept.ViaEdge != nil {
				o.M.NewTrack(
					o.RailNodes[deptID],          // from
					o.RailNodes[destID],          // to
					o.RailEdges[dept.ViaEdge.ID], // via
					dept.Value)                   // cost
			}
		}
	}
	o.ReRouting = false
	return payload.Route
}

func scanRail(o *entities.Player) *Model {
	model := NewModel()
	for _, rn := range o.RailNodes {
		model.AddGoalID(rn.ID)
		model.FindOrCreateNode(rn)
	}

	for _, re := range o.RailEdges {
		model.FindOrCreateEdge(re)
	}

	return model
}
