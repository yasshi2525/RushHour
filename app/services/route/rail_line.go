package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func SearchRailLine(l *entities.RailLine, parallel int) []*entities.Transport {
	results := []*entities.Transport{}
	template := scanRailLine(l)

	payload, _ := Search(context.Background(), entities.PLATFORM, parallel, template)

	for destID, model := range payload.Route {
		for deptID, dept := range model.Nodes[entities.PLATFORM] {
			if destID != deptID && dept.ViaEdge != nil {
				tr := &entities.Transport{
					l.Platforms[deptID],      // from
					l.Platforms[destID],      // to
					l.Tasks[dept.ViaEdge.ID], // via
					dept.Value}               // cost
				results = append(results, tr)
			} // ViaEdge = nil means cannot go to dest from dept by following line
		}
	}

	return results
}

func scanRailLine(l *entities.RailLine) *Model {
	model := NewModel()

	// gen goalid
	for _, p := range l.Platforms {
		model.AddGoalID(p.ID)
	}

	// gen nodes, edges
	for _, lt := range l.Tasks {
		model.FindOrCreateEdge(lt)
	}
	return model
}
