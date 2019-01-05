package entities

import (
	"testing"
)

func TestIsIn(t *testing.T) {
	scale := 4.0

	cases := []struct {
		target *Point
		want   bool
	}{
		{&Point{6, 6}, true},
		{&Point{-10, 0}, false},
		{&Point{10, 0}, false},
		{&Point{0, -10}, false},
		{&Point{0, 10}, false},
	}

	for _, c := range cases {
		if got := c.target.IsIn(0, 0, scale); got != c.want {
			t.Errorf("IsIn(%v) == %t, want %t", c.target, got, c.want)
		}
	}
}
