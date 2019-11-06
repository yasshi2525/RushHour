package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
)

func TestRailNode(t *testing.T) {
	t.Run("NewRailNode", func(t *testing.T) {
		a, _ := auth.GetAuther(config.CnfAuth{})
		m := NewModel(config.CnfEntity{}, a)
		o := m.NewPlayer()

		var x, y float64 = 10.0, 20.0

		rn := m.NewRailNode(o, x, y)

		TestCases{
			{"O", rn.O, o},
			{"X", rn.X, x},
			{"Y", rn.Y, y},
			{"o.rn", o.RailNodes[rn.Idx()], rn},
			{"model", m.RailNodes[rn.Idx()], rn},
			{"reroute", o.ReRouting, true},
		}.Assert(t)
	})

	t.Run("Extend", func(t *testing.T) {
		t.Run("without line", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			from := m.NewRailNode(o, 0, 0)
			var x, y float64 = 10.0, 20.0
			to, e1 := from.Extend(x, y)

			TestCases{
				{"to.X", to.X, x},
				{"to.Y", to.Y, y},
				{"e1.from", e1.FromNode, from},
				{"e1.to", e1.ToNode, to},
				{"e1.reverse", e1.Reverse.Reverse, e1},
				{"e2.from", e1.Reverse.FromNode, to},
				{"e2.to", e1.Reverse.ToNode, from},
				{"reroute", o.ReRouting, true},
			}.Assert(t)
		})
		t.Run("with autoExt", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			l.AutoExt = true

			n1 := m.NewRailNode(o, 0, 0)
			n2, e12 := n1.Extend(10, 0)
			head, _ := l.StartEdge(e12)

			_, e23 := n2.Extend(20, 0)

			TestCaseLineTasks{
				{"n1->n2", OnMoving, e12},
				{"n2->n3", OnMoving, e23},
				{"n3->n2", OnMoving, e23.Reverse},
				{"n2->n1", OnMoving, e12.Reverse},
				{"n1->n2", OnMoving, e12},
			}.Assert(t, head)

		})
	})

	t.Run("CheckDelete", func(t *testing.T) {
		t.Run("block overPlatform", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			st := m.NewStation(o)
			g := m.NewGate(st)
			rn := m.NewRailNode(o, 0, 0)
			m.NewPlatform(rn, g)

			if rn.CheckDelete() == nil {
				t.Error("wanted error, but got nil")
			}
		})
		t.Run("block lineTask", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			l := m.NewRailLine(o)

			n1 := m.NewRailNode(o, 0, 0)
			n2, e12 := n1.Extend(10, 0)
			_, e23 := n2.Extend(20, 0)

			_, tail := l.StartEdge(e12)
			tail.Stretch(e23)

			if n2.CheckDelete() == nil {
				t.Error("wanted error, but got nil")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("isolated", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			m.NewPlatform(rn, g)

			rn.CheckDelete()
			rn.Delete()
			if m.RailNodes[rn.ID] == rn {
				t.Error("rn remains in model")
			}
			TestCases{
				{"rn", len(m.RailNodes), 0},
				{"p", len(m.Platforms), 0},
				{"reroute", o.ReRouting, true},
				{"o", len(o.RailNodes), 0},
			}.Assert(t)
		})
		t.Run("line to isolated", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			n1 := m.NewRailNode(o, 0, 0)
			n2, _ := n1.Extend(10, 0)
			n2.Extend(20, 0)
			n2.CheckDelete()
			n2.Delete()
			TestCases{
				{"rn", len(m.RailNodes), 2},
				{"re", len(m.RailEdges), 0},
				{"reroute", o.ReRouting, true},
				{"o", len(o.RailNodes), 2},
			}.Assert(t)
		})
	})
}
