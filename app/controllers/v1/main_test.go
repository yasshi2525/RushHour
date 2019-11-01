package v1

import (
	"fmt"
	"math"
	"net/url"
	"testing"

	"github.com/yasshi2525/RushHour/app/services"
	"gopkg.in/go-playground/validator.v9"
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

func assertError(t *testing.T, in interface{}, want []string) {
	errs := buildErrorMessages(validate.Struct(in).(validator.ValidationErrors))
	if len(errs) != len(want) {
		t.Errorf("buildErrorMessages(%v) == %d errors, want %d errors", in, len(errs), len(want))
	} else {
		for i := 0; i < len(errs); i++ {
			got, want := errs[i], want[i]
			if fmt.Sprintf("%s", got) != want {
				t.Errorf("buildErrorMessages(%v)[%d] == %s, want %s", in, i, got, want)
			}
		}
	}
}

func TestBuildErrorMessages(t *testing.T) {
	Init()
	t.Run("GetGameMap", func(t *testing.T) {
		cases := []struct {
			in   gameMapRequest
			want []string
		}{
			{
				// too small Scale
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale-0.0001),
					Delegate: "0.0",
				},
				want: []string{"scale must be gte 0.000000"},
			},
			{
				// too large Scale
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MaxScale+0.0001),
					Delegate: "0.0",
				},
				want: []string{"scale must be lte 16.000000"},
			},
			{
				// too small Delegate
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: "-0.0001",
				},
				want: []string{"delegate must be gte 0"},
			},
			{
				// too large Delegate
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: fmt.Sprintf("%f", services.Config.Entity.MinScale+1),
				},
				want: []string{"delegate must be lte 0.000000"},
			},
			{
				// left over
				in: gameMapRequest{
					Cx:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cx must be gte -32767.500000"},
			},
			{
				// right over
				in: gameMapRequest{
					Cx:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cx must be lte 32767.500000"},
			},
			{
				// top over
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cy must be gte -32767.500000"},
			},
			{
				// bottom over
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
					Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cy must be lte 32767.500000"},
			},
			{
				// invalid format
				in: gameMapRequest{
					Cx:       "invalid",
					Cy:       "invalid",
					Scale:    "invalid",
					Delegate: "invalid",
				},
				want: []string{
					"cx must be numeric",
					"cy must be numeric",
					"scale must be numeric",
					"delegate must be numeric",
				},
			},
			{
				// empty
				in: gameMapRequest{},
				want: []string{
					"cx must be required",
					"cy must be required",
					"scale must be required",
					"delegate must be required",
				},
			},
		}
		for _, c := range cases {
			assertError(t, c.in, c.want)
		}
	})

	t.Run("Login", func(t *testing.T) {
		cases := []struct {
			in   loginRequest
			want []string
		}{
			{
				// invalid mail address
				in: loginRequest{
					ID:       "admin",
					Password: "password",
				},
				want: []string{"id must be email"},
			},
			{
				// empty
				in: loginRequest{},
				want: []string{
					"id must be required",
					"password must be required",
				},
			},
		}
		for _, c := range cases {
			assertError(t, c.in, c.want)
		}
	})

	t.Run("Register", func(t *testing.T) {
		cases := []struct {
			in   registerRequest
			want []string
		}{
			{
				// too small hue
				in: registerRequest{
					loginRequest: loginRequest{
						ID:       "test@example.com",
						Password: "password",
					},
					DisplayName: "",
					Hue:         "0",
				},
				want: []string{"hue must be gte 0"},
			},
			{
				// too large hue
				in: registerRequest{
					loginRequest: loginRequest{
						ID:       "test@example.com",
						Password: "password",
					},
					DisplayName: "",
					Hue:         "360",
				},
				want: []string{"hue must be lt 360"},
			},
		}
		for _, c := range cases {
			assertError(t, c.in, c.want)
		}
	})
}

func TestInit(t *testing.T) {
	Init()
	if validate == nil {
		t.Errorf("validate == nil, want not nil")
	}
}
