package entities

import (
	"testing"
	"time"
)

func TestNewNase(t *testing.T) {
	m := NewModel()
	t.Run("without owner", func(t *testing.T) {
		b := m.NewBase(COMPANY)
		if got := b.Idx(); got == 0 {
			t.Errorf("id should be > 0, but got = %d", got)
		}
		if got := b.M; got != m {
			t.Errorf("M should be %v, but got = %v", m, got)
		}
		if got := b.Type(); got != COMPANY {
			t.Errorf("T should be COMPANY, but got = %v", got)
		}
		if got := b.O; got != nil {
			t.Errorf("O should be nil, but got = %v", b.O)
		}
		if got := b.OwnerID; got != ZERO {
			t.Errorf("OwnerID should be 0, but got = %d", b.OwnerID)
		}
		if got := b.ChangedAt; got == *new(time.Time) {
			t.Errorf("ChangedAt should be > 0, but got = %v", got)
		}
	})
	t.Run("with owner", func(t *testing.T) {
		o := m.NewPlayer()
		b := m.NewBase(COMPANY, o)
		if b.O == nil {
			t.Errorf("O should not be nil, but got = nil")
		}
		if b.OwnerID == ZERO {
			t.Errorf("OwnerID should not be 0, but got = %d", ZERO)
		}
	})
}

func TestPermits(t *testing.T) {
	m := NewModel()
	my, oth, guest, admin := m.NewPlayer(), m.NewPlayer(), m.NewPlayer(), m.NewPlayer()
	my.Level, oth.Level, guest.Level, admin.Level = Normal, Normal, Guest, Admin

	cases := []struct {
		name   string
		b      Base
		target *Player
		want   bool
	}{
		{"Permits_Self", my.Base, my, true},
		{"Reject_Other", my.Base, oth, false},
		{"Reject_Guest", my.Base, guest, false},
		{"Permits_Admin", my.Base, admin, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.b.Permits(c.target); got != c.want {
				t.Errorf("want = %t, but got = %t", c.want, got)
			}
		})
	}
}

func TestIsChanged(t *testing.T) {
	m := NewModel()
	b := m.NewBase(RESIDENCE)

	cases := []struct {
		name  string
		b     Base
		after time.Time
		want  bool
	}{
		{"Past", b, time.Time{}, true},
		{"Future", b, time.Now(), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.b.IsChanged(c.after); got != c.want {
				t.Errorf("want = %t, but got = %t", c.want, got)
			}
		})
	}
}
