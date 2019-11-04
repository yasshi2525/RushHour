package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
)

func TestLineController(t *testing.T) {
	t.Run("Shrink", func(t *testing.T) {
		t.Run("head", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			_, e01 := n0.Extend(10, 0)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n0, g)

			l := m.NewRailLine(o)
			head := m.NewLineTaskDept(l, p)
			tail := m.NewLineTask(l, e01, head)

			head.Shrink(p)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
			}.Assert(t, tail)

			TestCases{
				{"tail.before", tail.before, (*LineTask)(nil)},
				{"tail.dept", tail.Dept, (*Platform)(nil)},
				{"model", len(m.LineTasks), 1},
			}.Assert(t)
		})

		t.Run("tail", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			l := m.NewRailLine(o)
			head := m.NewLineTask(l, e01)
			tail := m.NewLineTaskDept(l, p, head)

			tail.Shrink(p)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
			}.Assert(t, head)

			TestCases{
				{"head.before", head.next, (*LineTask)(nil)},
				{"head.dept", head.Dest, (*Platform)(nil)},
				{"head.next", head.next, (*LineTask)(nil)},
				{"model", len(m.LineTasks), 1},
			}.Assert(t)
		})
	})

	t.Run("Shave", func(t *testing.T) {
		t.Run("no reverse set next to nil", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			_, e12 := n1.Extend(20, 0)

			l := m.NewRailLine(o)
			head := m.NewLineTask(l, e01)
			target := m.NewLineTask(l, e12, head)

			target.Shave(e12)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
			}.Assert(t, head)

			TestCases{
				{"head.next", head.next, (*LineTask)(nil)},
				{"model", len(m.LineTasks), 1},
			}.Assert(t)
		})

		t.Run("no next set next to nil", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			_, e12 := n1.Extend(20, 0)

			l := m.NewRailLine(o)
			head := m.NewLineTask(l, e01)
			target := m.NewLineTask(l, e12, head)
			m.NewLineTask(l, e12.Reverse, target)

			target.Shave(e12)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
			}.Assert(t, head)

			TestCases{
				{"head.next", head.next, (*LineTask)(nil)},
				{"model", len(m.LineTasks), 1},
			}.Assert(t)
		})

		t.Run("delete redundant departure", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			_, e12 := n1.Extend(20, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			l := m.NewRailLine(o)
			head := m.NewLineTask(l, e01)
			tail := m.NewLineTaskDept(l, p, head)
			target := m.NewLineTask(l, e12, tail)
			tail = m.NewLineTask(l, e12.Reverse, target)
			tail = m.NewLineTaskDept(l, p, tail)
			tail = m.NewLineTask(l, e01.Reverse, tail)

			target.Shave(e12)

			TestCaseLineTasks{
				{"n0->n1", OnStopping, e01},
				{"n1", OnDeparture, p},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)

			TestCases{
				{"model", len(m.LineTasks), 3},
			}.Assert(t)
		})

		t.Run("change passing to stopping", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{
				MaxScale: 6,
			}, a)
			o := m.NewPlayer()
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			_, e12 := n1.Extend(20, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			l := m.NewRailLine(o)
			l.AutoPass = true
			head := m.NewLineTask(l, e01)
			target := m.NewLineTask(l, e12, head)
			l.AutoPass = false
			tail := m.NewLineTask(l, e12.Reverse, target)
			tail = m.NewLineTaskDept(l, p, tail)
			tail = m.NewLineTask(l, e01.Reverse, tail)

			target.Shave(e12)

			TestCaseLineTasks{
				{"n0->n1", OnStopping, e01},
				{"n1", OnDeparture, p},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)

			TestCases{
				{"model", len(m.LineTasks), 3},
			}.Assert(t)
		})
	})
}
