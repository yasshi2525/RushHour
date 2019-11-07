package route

import (
	"context"

	"github.com/yasshi2525/RushHour/entities"
)

// RefreshTransports set minimum distance route on specified rail line.
func RefreshTransports(l *entities.RailLine, parallel int) map[uint]*Model {
	l.ClearTransports()
	if !l.IsRing() || len(l.Trains) == 0 {
		return nil
	}
	template := scanRailLine(l)

	payload, _ := Search(context.Background(), entities.PLATFORM, parallel, template)

	for destID, model := range payload.Route {
		for deptID, dept := range model.Nodes[entities.PLATFORM] {
			if dept.ViaEdge != nil {
				l.M.NewTransport(
					l.Stops[deptID],          // from
					l.Stops[destID],          // to
					l.Tasks[dept.ViaEdge.ID], // via
					dept.Value)               // cost
			} // ViaEdge = nil means cannot go to dest from dept by following line
		}
	}
	l.ReRouting = false
	return payload.Route
}

func scanRailLine(l *entities.RailLine) *Model {
	model := NewModel()

	// gen goalid
	for _, p := range l.Stops {
		model.AddGoalID(p.ID)
		model.FindOrCreateNode(p)
	}

	// gen nodes, edges
	for _, lt := range l.Tasks {
		model.FindOrCreateEdge(lt)
	}
	return model
}
