package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func SearchRail(o *entities.Player, parallel int) (map[uint]*Model, []*entities.Track) {
	results := []*entities.Track{}
	template := scanRail(o)

	payload, _ := Search(context.Background(), entities.RAILNODE, parallel, template)

	for destID, model := range payload.Route {
		for deptID, dept := range model.Nodes[entities.RAILNODE] {
			if destID != deptID && dept.ViaEdge != nil {
				tr := &entities.Track{
					o.RailNodes[destID],          // from
					o.RailNodes[deptID],          // to
					o.RailEdges[dept.ViaEdge.ID], // via
					dept.Value}                   // cost
				results = append(results, tr)
			}
		}
	}
	return payload.Route, results
}

func scanRail(o *entities.Player) *Model {
	model := NewModel()

	for _, rn := range o.RailNodes {
		model.AddGoalID(rn.ID)
	}

	for _, re := range o.RailEdges {
		model.FindOrCreateEdge(re)
	}

	return model
}
