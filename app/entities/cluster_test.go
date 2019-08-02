package entities

import (
	"math"
	"testing"
)

func TestCluster(t *testing.T) {
	t.Run("NewCluster", func(t *testing.T) {
		Const = Config{MaxScale: 16}
		m := NewModel()

		root := m.NewCluster(nil, 0, 0)
		child := m.NewCluster(root, -1, 1)

		TestCases{
			{"root.scale", root.Scale, Const.MaxScale},
			{"child.scale", child.Scale, Const.MaxScale - 1},
			{"child.x", child.X, -math.Pow(2, Const.MaxScale) / 4},
			{"child.y", child.Y, math.Pow(2, Const.MaxScale) / 4},
		}.Assert(t)
	})

	t.Run("Init", func(t *testing.T) {
		m := NewModel()

		Const = Config{MaxScale: 2, MinScale: 1}

		parent := m.NewCluster(nil, 0, 0)
		child := m.NewCluster(parent, 0, 0)

		TestCases{
			{"parent.chPos", parent.ChPos[0][0] != nil, true},
			{"child.chPos", child.ChPos[0][0] == nil, true},
		}.Assert(t)
	})

	t.Run("FindChunk", func(t *testing.T) {
		t.Run("same", func(t *testing.T) {
			Const = Config{MaxScale: 16}
			m := NewModel()

			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)

			res := m.RootCluster.FindChunk(rn, Const.MaxScale)

			TestCases{
				{"data", res, m.RootCluster.Data[o.ID]},
			}.Assert(t)
		})

		t.Run("child", func(t *testing.T) {
			Const = Config{MaxScale: 2, MinScale: 1}
			m := NewModel()

			o := m.NewPlayer()
			rn := m.NewRailNode(o, -0.5, 0.5)

			res := m.RootCluster.FindChunk(rn, Const.MinScale)

			TestCases{
				{"data", res, m.RootCluster.Children[1][0].Data[o.ID]},
			}.Assert(t)
		})
	})
}
