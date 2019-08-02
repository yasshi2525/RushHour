package entities

import "testing"

func TestPlatform(t *testing.T) {
	Const = Config{MaxScale: 7}
	t.Run("NewPlatform", func(t *testing.T) {
		t.Run("stop", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			l.AutoExt = true
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			_, e12 := n1.Extend(20, 0)
			_, e13 := n1.Extend(10, 10)
			head, _ := l.StartEdge(e01)
			head.InsertRailEdge(e12)
			head.InsertRailEdge(e13)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			TestCases{
				{"O", p.O, o},
				{"rn", p.OnRailNode, n1},
				{"rnID", p.RailNodeID, n1.ID},
				{"g", p.WithGate, g},
				{"gID", p.GateID, g.ID},
				{"st", p.InStation, st},
				{"stID", p.StationID, st.ID},
				{"inSteps", len(p.InSteps()), 1},
				{"outSteps", len(p.OutSteps()), 1},
				{"inTasks", len(p.InTasks), 3},
				{"stayTasks", len(p.StayTasks), 3},
				{"outTasks", len(p.OutTasks), 3},
				{"model", m.Platforms[p.Idx()], p},
				{"s", len(m.Steps), 2},
			}.Assert(t)

			TestCaseLineTasks{
				{"n0->n1", OnStopping, e01},
				{"p", OnDeparture, p},
				{"n1->n3", OnMoving, e13},
				{"n3->n1", OnStopping, e13.Reverse},
				{"p", OnDeparture, p},
				{"n1->n2", OnMoving, e12},
				{"n2->n1", OnStopping, e12.Reverse},
				{"p", OnDeparture, p},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)
		})

		t.Run("pass", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			l.AutoExt = true
			l.AutoPass = true
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			head, _ := l.StartEdge(e01)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			TestCases{
				{"inTasks", len(p.InTasks), 1},
				{"stayTasks", len(p.StayTasks), 0},
				{"outTasks", len(p.OutTasks), 1},
			}.Assert(t)

			TestCaseLineTasks{
				{"n0->n1", OnPassing, e01},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("stop", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			l.AutoExt = true
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			head, _ := l.StartEdge(e01)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			p.Delete()

			TestCases{
				{"rn", n1.OverPlatform, (*Platform)(nil)},
				{"rnID", n1.PlatformID, uint(0)},
				{"o", len(o.Platforms), 0},
				{"s", len(m.Steps), 0},
				{"model", len(m.Platforms), 0},
			}.Assert(t)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)
		})

		t.Run("pass", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			l := m.NewRailLine(o)
			l.AutoExt = true
			l.AutoPass = true
			n0 := m.NewRailNode(o, 0, 0)
			n1, e01 := n0.Extend(10, 0)
			head, _ := l.StartEdge(e01)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n1, g)

			p.Delete()

			TestCases{
				{"rn", n1.OverPlatform, (*Platform)(nil)},
				{"rnID", n1.PlatformID, uint(0)},
				{"o", len(o.Platforms), 0},
				{"s", len(m.Steps), 0},
				{"model", len(m.Platforms), 0},
			}.Assert(t)

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)
		})

		t.Run("multiple", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			l := m.NewRailLine(o)

			n0 := m.NewRailNode(o, 0, 0)
			_, e01 := n0.Extend(10, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(n0, g)

			_, tail := l.StartPlatform(p)
			tail.InsertRailEdge(e01)
			tail.InsertRailEdge(e01)

			p.Delete()

			head, _ := l.Borders()

			TestCaseLineTasks{
				{"n0->n1", OnMoving, e01},
				{"n1->n0", OnMoving, e01.Reverse},
				{"n0->n1", OnMoving, e01},
				{"n1->n0", OnMoving, e01.Reverse},
			}.Assert(t, head)
		})
	})
}
