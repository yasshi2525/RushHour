package entities

import "testing"

func TestRailEdge(t *testing.T) {
	t.Run("NewRailEdge", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		n1 := m.NewRailNode(o, 0, 0)
		n2 := m.NewRailNode(o, 10, 0)
		re := m.NewRailEdge(n1, n2)

		TestCases{
			{"O", re.O, o},
			{"From", re.FromNode, n1},
			{"FromID", re.FromID, n1.ID},
			{"To", re.ToNode, n2},
			{"ToID", re.ToID, n2.ID},
			{"model", m.RailEdges[re.ID], re},
			{"reroute", o.ReRouting, true},
		}.Assert(t)
	})

}
