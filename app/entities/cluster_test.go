package entities

import (
	"math"
	"testing"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
)

func TestCluster(t *testing.T) {
	t.Run("NewCluster", func(t *testing.T) {
		a, _ := auth.GetAuther(config.CnfAuth{})
		m := NewModel(config.CnfEntity{
			MaxScale: 16,
		}, a)

		root := m.NewCluster(nil, 0, 0)
		child := m.NewCluster(root, -1, 1)

		TestCases{
			{"root.scale", root.Scale, m.conf.MaxScale},
			{"child.scale", child.Scale, m.conf.MaxScale - 1},
			{"child.x", child.X, -math.Pow(2, m.conf.MaxScale) / 4},
			{"child.y", child.Y, math.Pow(2, m.conf.MaxScale) / 4},
		}.Assert(t)
	})

	t.Run("Init", func(t *testing.T) {
		a, _ := auth.GetAuther(config.CnfAuth{})
		m := NewModel(config.CnfEntity{
			MaxScale: 2,
			MinScale: 1,
		}, a)

		parent := m.NewCluster(nil, 0, 0)
		child := m.NewCluster(parent, 0, 0)

		TestCases{
			{"parent.chPos", parent.ChPos[0][0] != nil, true},
			{"child.chPos", child.ChPos[0][0] == nil, true},
		}.Assert(t)
	})

	t.Run("FindChunk", func(t *testing.T) {
		t.Run("same", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
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
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 2,
				MinScale: 1,
			}, a)

			o := m.NewPlayer()
			rn := m.NewRailNode(o, -0.5, 0.5)

			res := m.RootCluster.FindChunk(rn, m.conf.MinScale)
			TestCases{
				{"data", res, m.RootCluster.Children[1][0].Data[o.ID]},
			}.Assert(t)
		})
	})
}
