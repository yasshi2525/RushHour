package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
)

func TestRailEdge(t *testing.T) {
	t.Run("NewRailEdge", func(t *testing.T) {
		a, _ := auth.GetAuther(config.CnfAuth{})
		m := NewModel(config.CnfEntity{}, a)
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
			{"from.Out", n1.OutEdges[re.ID], re},
			{"from.In", n1.InEdges[re.ReverseID], re.Reverse},
			{"to.In", n2.InEdges[re.ID], re},
			{"to.Out", n2.OutEdges[re.ReverseID], re.Reverse},
			{"o.re", o.RailEdges[re.ID], re},
			{"o.re.rev", o.RailEdges[re.ReverseID], re.Reverse},
			{"m.re", m.RailEdges[re.Idx()], re},
			{"m.re.rev", m.RailEdges[re.ReverseID], re.Reverse},
			{"reroute", o.ReRouting, true},
		}.Assert(t)
	})
	t.Run("CheckDelete", func(t *testing.T) {
		t.Run("relay", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			l := m.NewRailLine(o)

			n1 := m.NewRailNode(o, 0, 0)
			n2, e12 := n1.Extend(10, 0)
			_, e23 := n2.Extend(20, 0)

			_, tail := l.StartEdge(e12)
			tail.Stretch(e23)

			if e12.CheckDelete() == nil {
				t.Error("want err, but got nil")
			}
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("without LineTask", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, 0, 0)
			n2, re := n1.Extend(10, 0)
			re.CheckDelete()
			re.Delete()

			TestCases{
				{"n1.In", len(n1.InEdges), 0},
				{"n1.Out", len(n1.OutEdges), 0},
				{"n2.In", len(n2.InEdges), 0},
				{"n2.Out", len(n2.OutEdges), 0},
				{"o.re", len(o.RailEdges), 0},
				{"reroute", o.ReRouting, true},
				{"m.re", len(m.RailEdges), 0},
			}.Assert(t)
		})
		t.Run("with LineTask", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, 0, 0)
			_, re := n1.Extend(10, 0)
			l := m.NewRailLine(o)
			l.AutoExt = true
			l.StartEdge(re)

			re.CheckDelete()
			re.Delete()

			TestCases{
				{"re.lt", len(re.LineTasks), 0},
				{"m.lt", len(m.LineTasks), 0},
			}.Assert(t)
		})
	})
}
