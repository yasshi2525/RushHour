package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
)

func TestCluster(t *testing.T) {
	a, _ := auth.GetAuther(config.CnfAuth{Key: "----------------"})
	t.Run("NewCluster", func(t *testing.T) {
		m := NewModel(config.CnfEntity{
			MaxScale: 16,
		}, a)

		root := m.NewCluster(nil, 0, 0)
		child := m.NewCluster(root, 0, 1)

		TestCases{
			{"root.scale", root.Scale, m.conf.MaxScale},
			{"child.scale", child.Scale, m.conf.MaxScale - 1},
			{"child.x", child.X, 0},
			{"child.y", child.Y, 1},
		}.Assert(t)
	})

	t.Run("FindChunk", func(t *testing.T) {
		t.Run("same", func(t *testing.T) {
			m := NewModel(config.CnfEntity{
				MaxScale: 16,
			}, a)

			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)

			res := m.RootCluster.FindChunk(rn, m.conf.MaxScale)

			TestCases{
				{"data", res, m.RootCluster.Data[o.ID]},
			}.Assert(t)
		})

		t.Run("child", func(t *testing.T) {
			m := NewModel(config.CnfEntity{
				MaxScale: 2,
				MinScale: 1,
			}, a)

			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 2)

			res := m.RootCluster.FindChunk(rn, m.conf.MinScale)
			t.Logf("%+v", res)
			TestCases{
				{"data", res, m.RootCluster.Children[1][0].Data[o.ID]},
			}.Assert(t)
		})
	})
}
