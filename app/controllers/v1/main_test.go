package v1

import (
	"net/url"
	"testing"

	"github.com/yasshi2525/RushHour/app/services"
)

func TestMain(m *testing.M) {
	services.Config.Entity.MinScale = 0
	services.Config.Entity.MaxScale = 16
	m.Run()
}

func TestMapToStruct(t *testing.T) {
	// check certain map keyed "target"'s value insearting into struct
	type Sample struct {
		Target string `json:"target"`
	}

	cases := []struct {
		in struct {
			params url.Values
			out    *Sample
		}
		want string
	}{
		{
			in: struct {
				params url.Values
				out    *Sample
			}{
				params: map[string][]string{"target": []string{"true"}},
				out:    &Sample{},
			}, want: "true",
		}, {
			in: struct {
				params url.Values
				out    *Sample
			}{
				params: map[string][]string{"untarget": []string{"true"}},
				out:    &Sample{},
			}, want: "",
		},
	}
	for _, c := range cases {
		got := mapToStruct(c.in.params, c.in.out).(*Sample).Target
		if got != c.want {
			t.Errorf("mapToStruct(%v,%v).Target == %v, want %v", c.in.params, c.in.out, got, c.want)
		}
		if c.in.out.Target != c.want {
			t.Errorf("mapToStruct(%v,%v), 2nd params.Target == %v, want %v", c.in.params, c.in.out, c.in.out.Target, c.want)
		}
	}
}

func TestInit(t *testing.T) {
	Init()
	if validate == nil {
		t.Errorf("validate == nil, want not nil")
	}
}
