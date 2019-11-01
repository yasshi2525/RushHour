package v1

import (
	"fmt"
	"math"
	"testing"

	"github.com/yasshi2525/RushHour/app/services"
	"gopkg.in/go-playground/validator.v9"
)

func TestValidgameMapRequest(t *testing.T) {
	Init()

	cases := []struct {
		in   gameMapRequest
		want []string
	}{
		{
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: nil,
		}, {
			// too small Scale
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale-0.0001),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'gte' tag"},
		}, {
			// too large Scale
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MaxScale+0.0001),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'lte' tag"},
		}, {
			// too small Delegate
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "-0.0001",
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'gte' tag"},
		}, {
			// too large Delegate
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: fmt.Sprintf("%f", services.Config.Entity.MinScale+1),
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'lte' tag"},
		}, {
			// left over
			in: gameMapRequest{
				Cx:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'gte' tag"},
		}, {
			// right over
			in: gameMapRequest{
				Cx:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'lte' tag"},
		}, {
			// top over
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", -math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cy' Error:Field validation for 'cy' failed on the 'gte' tag"},
		}, {
			// bottom over
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", math.Pow(2, services.Config.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", services.Config.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cy' Error:Field validation for 'cy' failed on the 'lte' tag"},
		}, {
			// invalid format
			in: gameMapRequest{
				Cx:       "invalid",
				Cy:       "invalid",
				Scale:    "invalid",
				Delegate: "invalid",
			},
			want: []string{
				"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.cy' Error:Field validation for 'cy' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'numeric' tag",
			},
		}, {
			// empty
			in: gameMapRequest{},
			want: []string{
				"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'required' tag",
				"Key: 'gameMapRequest.cy' Error:Field validation for 'cy' failed on the 'required' tag",
				"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'required' tag",
				"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'required' tag",
			},
		},
	}
	for _, c := range cases {
		if rawResult := validate.Struct(c.in); rawResult == nil {
			if c.want != nil {
				t.Errorf("validgameMapRequest(%v) == nil, want %v", c.in, c.want)
			}
		} else {
			res := rawResult.(validator.ValidationErrors)

			if len(res) != len(c.want) {
				t.Errorf("validgameMapRequest(%v) == %d errors, want %d errors", c.in, len(res), len(c.want))
			} else {
				for i := 0; i < len(res); i++ {
					got, want := res[i], c.want[i]
					if fmt.Sprintf("%s", got) != want {
						t.Errorf("validgameMapRequest(%v)[%d] == %s, want %s", c.in, i, got, want)
					}
				}
			}
		}
	}
}
