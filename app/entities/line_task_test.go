package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
)

func TestLineTask(t *testing.T) {
	t.Run("NewLineTaskDept", func(t *testing.T) {
		t.Run("new", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(rn, g)
			l := m.NewRailLine(o)

			lt := m.NewLineTaskDept(l, p)

			TestCases{
				{"O", lt.O, o},
				{"O.lt", o.LineTasks[lt.Idx()], lt},
				{"l.tasks", l.Tasks[lt.Idx()], lt},
				{"l.p", l.Stops[p.Idx()], p},
				{"l.reroute", l.ReRouting, true},
				{"lt.task", lt.TaskType, OnDeparture},
				{"lt.stay", lt.Stay, p},
				{"lt.dept", lt.Dept, p},
				{"lt.dest", lt.Dest, p},
				{"p.stay", p.StayTasks[lt.Idx()], lt},
				{"rn.out", rn.OutTasks[lt.Idx()], lt},
				{"rn.in", rn.InTasks[lt.Idx()], lt},
				{"model", m.LineTasks[lt.Idx()], lt},
			}.Assert(t)
		})

		t.Run("extend", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			l := m.NewRailLine(o)
			lt := m.NewLineTask(l, e01)
			tail := m.NewLineTaskDept(l, p, lt)

			TestCases{
				{"lt.next", lt.next, tail},
				{"tail.before", tail.before, lt},
			}.Assert(t)
		})
	})

	t.Run("NewLineTask", func(t *testing.T) {
		t.Run("new", func(t *testing.T) {
			t.Run("moving", func(t *testing.T) {
				a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
				o := m.NewPlayer()
				n0 := m.NewRailNode(o, 0, 0)
				n1, e01 := n0.Extend(10, 0)

				l := m.NewRailLine(o)
				lt := m.NewLineTask(l, e01)

				TestCases{
					{"O", lt.O, o},
					{"O.lt", o.LineTasks[lt.Idx()], lt},
					{"l.tasks", l.Tasks[lt.Idx()], lt},
					{"l.re", l.RailEdges[e01.Idx()], e01},
					{"l.reroute", l.ReRouting, true},
					{"lt.task", lt.TaskType, OnMoving},
					{"lt.moving", lt.Moving, e01},
					{"n0.out", n0.OutTasks[lt.Idx()], lt},
					{"n1.in", n1.InTasks[lt.Idx()], lt},
					{"re.tasks", e01.LineTasks[lt.Idx()], lt},
					{"model", m.LineTasks[lt.Idx()], lt},
				}.Assert(t)
			})

			t.Run("stopping", func(t *testing.T) {
				a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
				o := m.NewPlayer()
				n0 := m.NewRailNode(o, 0, 0)
				n1, e01 := n0.Extend(10, 0)
				st := m.NewStation(o)
				g := m.NewGate(st)
				p := m.NewPlatform(n1, g)

				l := m.NewRailLine(o)
				lt := m.NewLineTask(l, e01)
				TestCases{
					{"lt.task", lt.TaskType, OnStopping},
					{"lt.dest", lt.Dest, p},
					{"p.in", p.InTasks[lt.Idx()], lt},
				}.Assert(t)
			})

			t.Run("passing", func(t *testing.T) {
				a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
				o := m.NewPlayer()
				n0 := m.NewRailNode(o, 0, 0)
				n1, e01 := n0.Extend(10, 0)
				st := m.NewStation(o)
				g := m.NewGate(st)
				m.NewPlatform(n1, g)

				l := m.NewRailLine(o)
				l.AutoPass = true
				lt := m.NewLineTask(l, e01)
				TestCases{
					{"lt.task", lt.TaskType, OnPassing},
				}.Assert(t)
			})
		})

		t.Run("extend", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n0, g)
			_, e01 := n0.Extend(10, 0)

			l := m.NewRailLine(o)
			lt := m.NewLineTaskDept(l, p)
			tail := m.NewLineTask(l, e01, lt)

			TestCases{
				{"lt.next", lt.next, tail},
				{"tail.before", tail.before, lt},
			}.Assert(t)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("dept", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			rn := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(rn, g)

			lt := m.NewLineTaskDept(l, p)
			lt.Delete()

			TestCases{
				{"o.lt", len(o.LineTasks), 0},
				{"p.lt", len(p.StayTasks), 0},
				{"rn.in", len(rn.InTasks), 0},
				{"rn.out", len(rn.OutTasks), 0},
				{"l.tasks", len(l.Tasks), 0},
			}.Assert(t)
		})

		t.Run("moving", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			lt := m.NewLineTask(l, e01)
			lt.Delete()

			TestCases{
				{"o.lt", len(o.LineTasks), 0},
				{"p.lt", len(p.StayTasks), 0},
				{"n0.out", len(n0.OutTasks), 0},
				{"n1.in", len(n1.InTasks), 0},
				{"n1.out", len(n1.OutTasks), 0},
				{"l.tasks", len(l.Tasks), 0},
			}.Assert(t)
		})

	})
}
