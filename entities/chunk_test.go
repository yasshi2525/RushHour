package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
)

func TestChunk(t *testing.T) {
	a, _ := auth.GetAuther(config.CnfAuth{Key: "----------------"})
	t.Run("Add", func(t *testing.T) {
		t.Run("RailNode", func(t *testing.T) {
			m := NewModel(config.CnfEntity{
				MaxScale: 1,
			}, a)
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, 0, 0)
			n2 := m.NewRailNode(o, 0, 0)

			ch := m.RootCluster.Data[o.ID]

			TestCases{
				{"RailNode", ch.RailNode != nil, true},
				{"n1", ch.RailNode.List[n1.ID], n1},
				{"n2", ch.RailNode.List[n2.ID], n2},
			}.Assert(t)
		})

		t.Run("RailEdge", func(t *testing.T) {
			m := NewModel(config.CnfEntity{
				MaxScale: 2,
			}, a)
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, -0.5, -0.5)
			n2, re := n1.Extend(0.5, 0.5)

			from := m.RootCluster.FindChunk(n1, m.conf.MinScale)
			to := m.RootCluster.FindChunk(n2, m.conf.MinScale)
			dreFrom := from.OutRailEdges[re.ID]
			dreTo := to.InRailEdges[re.ID]

			TestCases{
				{"re", dreFrom == dreTo, true},
			}.Assert(t)
		})
	})
}
