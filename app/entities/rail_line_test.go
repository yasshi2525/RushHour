package entities

import "testing"

func TestRailLine(t *testing.T) {
	t.Run("NewRailLine", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		l := m.NewRailLine(o)

		TestCases{
			{"O", l.O, o},
			{"o.l", o.RailLines[l.Idx()], l},
			{"model", m.RailLines[l.Idx()], l},
		}.Assert(t)
	})

	t.Run("StartPlatform", func(t *testing.T) {
		t.Run("auto ext", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			from := m.NewRailNode(o, 0, 0)
			_, re := from.Extend(10, 0)

			m.NewTrack(from, from, re, 1)

			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(from, g)
			l := m.NewRailLine(o)
			l.AutoExt = true

			head, _ := l.StartPlatform(p)

			TestCaseLineTasks{
				{"n0", OnDeparture, p},
				{"n0->n1", OnMoving, re},
				{"n1->n0", OnStopping, re.Reverse},
			}.Assert(t, head)
		})
		t.Run("manual", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(rn, g)
			l := m.NewRailLine(o)

			head, _ := l.StartPlatform(p)

			TestCaseLineTasks{
				{"dep", OnDeparture, p},
			}.Assert(t, head)
		})
		t.Run("auto pass", func(t *testing.T) {
			m := NewModel()
			o := m.NewPlayer()
			rn := m.NewRailNode(o, 0, 0)
			st := m.NewStation(o)
			g := m.NewGate(st)
			p := m.NewPlatform(rn, g)
			l := m.NewRailLine(o)
			l.AutoPass = true

			head, tail := l.StartPlatform(p)

			TestCases{
				{"head", head, (*LineTask)(nil)},
				{"tail", tail, (*LineTask)(nil)},
				{"l", len(l.Tasks), 0},
			}.Assert(t)
		})
	})

	t.Run("Complement", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		from := m.NewRailNode(o, 0, 0)
		to, re := from.Extend(10, 0)

		m.NewTrack(from, from, re, 2)
		m.NewTrack(to, from, re.Reverse, 1)

		st := m.NewStation(o)
		g := m.NewGate(st)
		p := m.NewPlatform(from, g)
		l := m.NewRailLine(o)

		head, _ := l.StartPlatform(p)
		l.Complement()

		TestCaseLineTasks{
			{"n0", OnDeparture, p},
			{"n0->n1", OnMoving, re},
			{"n1->n0", OnStopping, re.Reverse},
		}.Assert(t, head)
	})

	t.Run("Delete", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		from := m.NewRailNode(o, 0, 0)
		_, re := from.Extend(10, 0)
		l := m.NewRailLine(o)
		l.AutoExt = true
		l.StartEdge(re)

		l.Delete()

		TestCases{
			{"lt", len(m.LineTasks), 0},
			{"model", len(m.RailLines), 0},
		}.Assert(t)
	})
}
