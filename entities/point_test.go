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

func TestLogarithm(t *testing.T) {
	type in struct {
		num   float64
		scale int
	}

	cases := []struct {
		in   in
		want int
	}{
		{in{num: 0.0, scale: 0}, 0},
		{in{num: 1.0, scale: 0}, 1},
		{in{num: 1.5, scale: 0}, 1},
		{in{num: 2.0, scale: 0}, 2},
		{in{num: 1.0, scale: 1}, 0},
		{in{num: 2.0, scale: 1}, 1},
		{in{num: 2.1, scale: 1}, 1},
		{in{num: 4.0, scale: 1}, 2},
		{in{num: 3.0, scale: 2}, 0},
		{in{num: 4.0, scale: 2}, 1},
		{in{num: 4.1, scale: 2}, 1},
		{in{num: 0.0, scale: -1}, 0},
		{in{num: 0.4, scale: -1}, 0},
		{in{num: 0.5, scale: -1}, 1},
		{in{num: 0.6, scale: -1}, 1},
		{in{num: 0.24, scale: -2}, 0},
		{in{num: 0.25, scale: -2}, 1},
		{in{num: 0.26, scale: -2}, 1},
	}

	for _, c := range cases {
		if got := Logarithm(c.in.num, c.in.scale); got != c.want {
			t.Errorf("Logarithm(%.2f,%d) got %d, want %d",
				c.in.num, c.in.scale, got, c.want)
		}
	}
}

func TestDeLogarithm(t *testing.T) {
	type in struct {
		num   int
		scale int
	}

	cases := []struct {
		in   in
		want float64
	}{
		{in{num: 0, scale: 0}, 0.0},
		{in{num: 1, scale: 0}, 1.0},
		{in{num: 2, scale: 0}, 2.0},
		{in{num: 0, scale: 1}, 0.0},
		{in{num: 1, scale: 1}, 2.0},
		{in{num: 2, scale: 1}, 4.0},
		{in{num: 3, scale: 1}, 6.0},
		{in{num: 0, scale: 2}, 0.0},
		{in{num: 1, scale: 2}, 4.0},
		{in{num: 2, scale: 2}, 8.0},
		{in{num: 0, scale: -1}, 0.0},
		{in{num: 1, scale: -1}, 0.5},
		{in{num: 2, scale: -1}, 1.0},
		{in{num: 3, scale: -1}, 1.5},
		{in{num: 0, scale: -2}, 0.00},
		{in{num: 1, scale: -2}, 0.25},
		{in{num: 2, scale: -2}, 0.50},
	}

	for _, c := range cases {
		if got := DeLogarithm(c.in.num, c.in.scale); got != c.want {
			t.Errorf("DeLogarithm(%d,%d) got %.2f, want %.2f",
				c.in.num, c.in.scale, got, c.want)
		}
	}
}
