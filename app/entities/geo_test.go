package entities

import (
	"testing"
)

func TestDist(t *testing.T) {
	pivot := &Point{0, 0}

	cases := []struct {
		in   *Point
		want float64
	}{
		{&Point{0, 0}, 0},
		{&Point{3, 4}, 5},
	}

	for _, c := range cases {
		if got := pivot.Dist(c.in); got != c.want {
			t.Errorf("Dist(%v) == %f, want %f", c.in, got, c.want)
		}
	}
}

func TestIsIn(t *testing.T) {
	center := &Point{0, 0}
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
		if got := c.target.IsIn(center, scale); got != c.want {
			t.Errorf("IsIn(%v) == %t, want %t", c.target, got, c.want)
		}
	}
}

func TestCost(t *testing.T) {
	from := NewResidence(1, 0, 0)
	to := NewCompany(1, 3, 4)

	step := &Step{1, from, to, 2.0}

	if got := step.Cost(); got != 10.0 {
		t.Errorf("Cost(%v, %v) == %f, want %f", from, to, got, 10.0)
	}
}
