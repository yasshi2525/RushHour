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
		m := NewModel()
		o := m.NewPlayer()
		rn := m.RailNodes(o, 0, 0)
		st := m.NewStation(st)

		p := m.NewPlatform(rn)
		l := m.NewRailLine(o)
	})
}
