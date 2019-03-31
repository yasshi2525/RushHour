package entities

import (
	"testing"
	"time"
)

func TestBase(t *testing.T) {
	t.Run("NewBase", func(t *testing.T) {
		m := NewModel()
		t.Run("without owner", func(t *testing.T) {
			b := m.NewBase(COMPANY)

			TestCases{
				{"id", b.Idx(), uint(1)},
				{"M", b.M, m},
				{"Type", b.Type(), COMPANY},
				{"O", b.O, (*Player)(nil)},
				{"OwnerID", b.OwnerID, uint(0)},
			}.Assert(t)

			if got := b.ChangedAt; got == (time.Time{}) {
				t.Errorf("ChangedAt should be > 0, but got = %v", got)
			}
		})
		t.Run("with owner", func(t *testing.T) {
			o := m.NewPlayer()
			b := m.NewBase(COMPANY, o)

			TestCases{
				{"O", b.O, o},
				{"OwnerID", b.OwnerID, o.Idx()},
			}.Assert(t)
		})
	})

	t.Run("Permits", func(t *testing.T) {
		m := NewModel()
		my, oth, guest, admin := m.NewPlayer(), m.NewPlayer(), m.NewPlayer(), m.NewPlayer()
		my.Level, oth.Level, guest.Level, admin.Level = Normal, Normal, Guest, Admin

		TestCases{
			{"Permits_Self", my, true},
			{"Reject_Other", oth, false},
			{"Reject_Guest", guest, false},
			{"Permits_Admin", admin, true},
		}.Assert(t, func(val interface{}) interface{} {
			return my.Permits(val.(*Player))
		})
	})

	t.Run("IsChanged", func(t *testing.T) {
		m := NewModel()
		b := m.NewBase(RESIDENCE)

		TestCases{
			{"Past", time.Time{}, true},
			{"Future", time.Now(), false},
		}.Assert(t, func(val interface{}) interface{} {
			return b.IsChanged(val.(time.Time))
		})
	})

}
