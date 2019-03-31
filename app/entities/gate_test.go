package entities

import "testing"

func TestGate(t *testing.T) {
	t.Run("NewGate", func(t *testing.T) {
		m := NewModel()
		m.NewCompany(0, 0)
		m.NewResidence(0, 0)
		o := m.NewPlayer()
		st := m.NewStation(o)
		g := m.NewGate(st)

		TestCases{
			{"O", g.O, o},
			{"st", g.InStation, st},
			{"stID", g.StationID, st.ID},
			{"st.g", st.Gate, g},
			{"st.gID", st.GateID, g.ID},
			{"model", m.Gates[g.Idx()], g},
			{"s", len(m.Steps), 3},
		}.Assert(t)
	})
	t.Run("Delete", func(t *testing.T) {
		m := NewModel()
		m.NewCompany(0, 0)
		m.NewResidence(0, 0)
		o := m.NewPlayer()
		st := m.NewStation(o)
		g := m.NewGate(st)

		g.CheckDelete()
		g.Delete()

		TestCases{
			{"o", len(o.Gates), 0},
			{"s", len(m.Steps), 1},
			{"model", len(m.Gates), 0},
		}.Assert(t)
	})
}
