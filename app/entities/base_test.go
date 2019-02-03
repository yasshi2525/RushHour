package entities

import (
	"testing"
	"time"
)

func TestPermits(t *testing.T) {
	InitType()
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
	InitType()
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
