package entities

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
)

func TestStation(t *testing.T) {
	t.Run("NewStation", func(t *testing.T) {
		a, _ := auth.GetAuther(config.CnfAuth{})
		m := NewModel(config.CnfEntity{}, a)
		o := m.NewPlayer()
		st := m.NewStation(o)

		TestCases{
			{"O", st.O, o},
			{"model", m.Stations[st.Idx()], st},
		}.Assert(t)
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("isolate", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			st := m.NewStation(o)

			st.CheckDelete()
			st.Delete()

			TestCases{
				{"o", len(o.Stations), 0},
				{"model", len(m.Stations), 0},
			}.Assert(t)
		})
		t.Run("station", func(t *testing.T) {
			a, _ := auth.GetAuther(config.CnfAuth{})
			m := NewModel(config.CnfEntity{}, a)
			o := m.NewPlayer()
			st := m.NewStation(o)
			g := m.NewGate(st)
			rn := m.NewRailNode(o, 0, 0)
			m.NewPlatform(rn, g)

			st.CheckDelete()
			st.Delete()

			TestCases{
				{"o.g", len(o.Gates), 0},
				{"o.p", len(o.Platforms), 0},
				{"m.g", len(m.Gates), 0},
				{"m.p", len(m.Platforms), 0},
			}.Assert(t)
		})
	})
}
