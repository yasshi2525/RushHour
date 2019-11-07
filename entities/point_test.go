package entities

import (
	"testing"
)

func TestIsIn(t *testing.T) {
	scale := 4.0

	cases := []struct {
		name string
		p    *Point
		want bool
	}{
		{"OK", &Point{6, 6}, true},
		{"TooRight", &Point{-10, 0}, false},
		{"TooLeft", &Point{10, 0}, false},
		{"TooUp", &Point{0, -10}, false},
		{"TooDown", &Point{0, 10}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.p.IsIn(0, 0, scale); got != c.want {
				t.Errorf("IsIn(%v) == %t, want %t", c.p, got, c.want)
			}
		})
	}
}
