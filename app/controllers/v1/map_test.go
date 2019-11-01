package v1

import (
	"fmt"
	"math"
	"testing"

	"github.com/yasshi2525/RushHour/app/services"
	"gopkg.in/go-playground/validator.v9"
)

func TestValidGameMapRequest(t *testing.T) {
	Init()

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
			want: []string{"Key: 'GameMapRequest.scale' Error:Field validation for 'scale' failed on the 'gte' tag"},
		},
		{
			// too large Scale
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MaxScale+0.0001),
				Delegate: "0.0",
			},
			want: []string{"Key: 'GameMapRequest.scale' Error:Field validation for 'scale' failed on the 'lte' tag"},
		},
		{
			// too small Delegate
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "-0.0001",
			},
			want: []string{"Key: 'GameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'gte' tag"},
		},
		{
			// too large Delegate
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: fmt.Sprintf("%f", services.Config.Entity.MinScale+1),
			},
			want: []string{"Key: 'GameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'lte' tag"},
		},
		{
			// left over
			in: GameMapRequest{
				Cx:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'GameMapRequest.cx' Error:Field validation for 'cx' failed on the 'gte' tag"},
		},
		{
			// right over
			in: GameMapRequest{
				Cx:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'GameMapRequest.cx' Error:Field validation for 'cx' failed on the 'lte' tag"},
		},
		{
			// top over
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'GameMapRequest.cy' Error:Field validation for 'cy' failed on the 'gte' tag"},
		},
		{
			// bottom over
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'GameMapRequest.cy' Error:Field validation for 'cy' failed on the 'lte' tag"},
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
				"Key: 'GameMapRequest.cx' Error:Field validation for 'cx' failed on the 'numeric' tag",
				"Key: 'GameMapRequest.cy' Error:Field validation for 'cy' failed on the 'numeric' tag",
				"Key: 'GameMapRequest.scale' Error:Field validation for 'scale' failed on the 'numeric' tag",
				"Key: 'GameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'numeric' tag",
			},
		},
		{
			// empty
			in: GameMapRequest{},
			want: []string{
				"Key: 'GameMapRequest.cx' Error:Field validation for 'cx' failed on the 'required' tag",
				"Key: 'GameMapRequest.cy' Error:Field validation for 'cy' failed on the 'required' tag",
				"Key: 'GameMapRequest.scale' Error:Field validation for 'scale' failed on the 'required' tag",
				"Key: 'GameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'required' tag",
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
}

func TestErrorGetGameMap(t *testing.T) {
	Init()
	cases := []struct {
		in   GameMapRequest
		want []string
	}{
		{
			// too small Scale
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale-0.0001),
				Delegate: "0.0",
			},
			want: []string{"scale must be gte 0.000000"},
		},
		{
			// too large Scale
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MaxScale+0.0001),
				Delegate: "0.0",
			},
			want: []string{"scale must be lte 16.000000"},
		},
		{
			// too small Delegate
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "-0.0001",
			},
			want: []string{"delegate must be gte 0"},
		},
		{
			// too large Delegate
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: fmt.Sprintf("%f", services.Config.Entity.MinScale+1),
			},
			want: []string{"delegate must be lte 0.000000"},
		},
		{
			// left over
			in: GameMapRequest{
				Cx:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"cx must be gte -32767.500000"},
		},
		{
			// right over
			in: GameMapRequest{
				Cx:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"cx must be lte 32767.500000"},
		},
		{
			// top over
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"cy must be gte -32767.500000"},
		},
		{
			// bottom over
			in: GameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"cy must be lte 32767.500000"},
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
				"cx must be numeric",
				"cy must be numeric",
				"scale must be numeric",
				"delegate must be numeric",
			},
		},
		{
			// empty
			in: GameMapRequest{},
			want: []string{
				"cx must be required",
				"cy must be required",
				"scale must be required",
				"delegate must be required",
			},
		},
	}
	for _, c := range cases {
		errs := errorGetGameMap(validate.Struct(c.in).(validator.ValidationErrors))

		if len(errs) != len(c.want) {
			t.Errorf("errorGetGameMap(%v) == %d errors, want %d errors", c.in, len(errs), len(c.want))
		} else {
			for i := 0; i < len(errs); i++ {
				got, want := errs[i], c.want[i]
				if fmt.Sprintf("%s", got) != want {
					t.Errorf("errorGetGameMap(%v)[%d] == %s, want %s", c.in, i, got, want)
				}
			}
		}
	}
}
