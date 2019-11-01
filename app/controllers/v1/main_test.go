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

func TestInit(t *testing.T) {
	Init()
	if validate == nil {
		t.Errorf("validate == nil, want not nil")
	} else {
		t.Run("validGameMapRequest", func(t *testing.T) {
			cases := []struct {
				in   GameMapRequest
				want []string
			}{
				{
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "0.0",
					},
					want: nil,
				},
				{
					// too small Scale
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale-0.0001),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Scale' Error:Field validation for 'Scale' failed on the 'gte' tag"},
				},
				{
					// too large Scale
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MaxScale+0.0001),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Scale' Error:Field validation for 'Scale' failed on the 'lte' tag"},
				},
				{
					// too small Delegate
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "-0.0001",
					},
					want: []string{"Key: 'GameMapRequest.Delegate' Error:Field validation for 'Delegate' failed on the 'gte' tag"},
				},
				{
					// too large Delegate
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: fmt.Sprintf("%f", services.Config.Entity.MinScale+1),
					},
					want: []string{"Key: 'GameMapRequest.Delegate' Error:Field validation for 'Delegate' failed on the 'lte' tag"},
				},
				{
					// left over
					in: GameMapRequest{
						Cx:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Cx' Error:Field validation for 'Cx' failed on the 'gte' tag"},
				},
				{
					// right over
					in: GameMapRequest{
						Cx:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
						Cy:       "0.0",
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Cx' Error:Field validation for 'Cx' failed on the 'lte' tag"},
				},
				{
					// top over
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Cy' Error:Field validation for 'Cy' failed on the 'gte' tag"},
				},
				{
					// bottom over
					in: GameMapRequest{
						Cx:       "0.0",
						Cy:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
						Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
						Delegate: "0.0",
					},
					want: []string{"Key: 'GameMapRequest.Cy' Error:Field validation for 'Cy' failed on the 'lte' tag"},
				},
				{
					// invalid format
					in: GameMapRequest{
						Cx:       "invalid",
						Cy:       "invalid",
						Scale:    "invalid",
						Delegate: "invalid",
					},
					want: []string{
						"Key: 'GameMapRequest.Cx' Error:Field validation for 'Cx' failed on the 'numeric' tag",
						"Key: 'GameMapRequest.Cy' Error:Field validation for 'Cy' failed on the 'numeric' tag",
						"Key: 'GameMapRequest.Scale' Error:Field validation for 'Scale' failed on the 'numeric' tag",
						"Key: 'GameMapRequest.Delegate' Error:Field validation for 'Delegate' failed on the 'numeric' tag",
					},
				},
			}
			for _, c := range cases {

				if rawResult := validate.Struct(c.in); rawResult == nil {
					if c.want != nil {
						t.Errorf("validGameMapRequest(%v) == nil, want %v", c.in, c.want)
					}
				} else {
					res := rawResult.(validator.ValidationErrors)

					if len(res) != len(c.want) {
						t.Errorf("validGameMapRequest(%v) == %d errors, want %d errors", c.in, len(res), len(c.want))
					} else {
						for i := 0; i < len(res); i++ {
							got, want := res[i], c.want[i]
							if fmt.Sprintf("%s", got) != want {
								t.Errorf("validGameMapRequest(%v)[%d] == %s, want %s", c.in, i, got, want)
							}
						}
					}
				}
			}
		})
	}
}
