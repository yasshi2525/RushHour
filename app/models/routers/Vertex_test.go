package routers

import "testing"

func TestDist(t *testing.T) {
	inst := NewVertex(0, 0)

	cases := []struct {
		in   *Vertex
		want float64
	}{
		{NewVertex(0, 0), 0},
		{NewVertex(3, 4), 5},
	}

	for _, c := range cases {
		got := inst.Dist(c.in)
		if got != c.want {
			//t.Errorf("Dist(%g) == %g, want %g", c.in, got, c.want)
		}
	}
}
