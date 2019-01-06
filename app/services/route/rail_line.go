package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func SearchRailLine(l *entities.RailLine, parallel int) []*entities.Track {
	results := []*entities.Track{}
	template := scanRailLine(l)

	payload, _ := Search(context.Background(), entities.PLATFORM, parallel, template)

	for _, dept := range l.Platforms {
		for destID, model := range payload.Route {
			// prevent self-relation
			if dept.ID == destID {
				continue
			}
			for _, n := range model.Nodes {
				if n.SameAs(dept) {
					if n.ViaEdge != nil {
						tr := entities.NewTrack(
							dept,                  // from
							l.Platforms[destID],   // to
							l.Tasks[n.ViaEdge.ID], // via
							n.Value)               // cost
						results = append(results, tr)
					} // ViaEdge = nil means cannot go to dest from dept by following line
					break
				}
			}
		}
	}
	return results
}

func scanRailLine(l *entities.RailLine) *Model {
	// gen pids
	pids := make([]uint, len(l.Platforms))
	i := 0
	for _, p := range l.Platforms {
		pids[i] = p.ID
		i++
	}

	// gen nodes, edges
	ns := NodeQueue{}
	es := make([]*Edge, len(l.Tasks))
	i = 0
	for _, lt := range l.Tasks {
		n1 := ns.AppendIfNotExists(lt.From())
		n2 := ns.AppendIfNotExists(lt.To())
		es[i] = NewEdge(lt, n1, n2)
		i++
	}
	return &Model{pids, ns, es}
}
