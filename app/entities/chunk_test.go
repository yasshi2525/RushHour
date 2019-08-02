package entities

import "testing"

func TestChunk(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		Const = Config{MaxScale: 1, MinScale: 0}

		t.Run("RailNode", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, 0, 0)
			n2 := m.NewRailNode(o, 0, 0)

			ch := m.RootCluster.Data[o.ID]

			TestCases{
				{"RailNode", ch.RailNode != nil, true},
				{"n1", ch.RailNode.RailNodes[n1.ID], n1},
				{"n2", ch.RailNode.RailNodes[n2.ID], n2},
			}.Assert(t)
		})

		t.Run("RailEdge", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, -0.5, -0.5)
			n2, re := n1.Extend(0.5, 0.5)

			from := m.RootCluster.FindChunk(n1, Const.MinScale)
			to := m.RootCluster.FindChunk(n2, Const.MinScale)
			dreFrom := from.OutRailEdges[re.ID]
			dreTo := to.InRailEdges[re.ID]

			TestCases{
				{"re", dreFrom == dreTo, true},
			}.Assert(t)
		})
	})
}
