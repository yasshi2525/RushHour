package v1

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

func TestValidgameMapRequest(t *testing.T) {
	v := initValidate()
	cases := []struct {
		in   gameMapRequest
		want []string
	}{
		{
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: "0.0",
			},
			want: nil,
		}, {
			// too small Scale
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale-0.0001),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'gte' tag"},
		}, {
			// too large Scale
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MaxScale+0.0001),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'lte' tag"},
		}, {
			// too small Delegate
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: "-0.0001",
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'gte' tag"},
		}, {
			// too large Delegate
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: fmt.Sprintf("%f", conf.Game.Entity.MinScale+1),
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'lte' tag"},
		}, {
			// left over
			in: gameMapRequest{
				Cx:       fmt.Sprintf("%f", -math.Pow(2, conf.Game.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'gte' tag"},
		}, {
			// right over
			in: gameMapRequest{
				Cx:       fmt.Sprintf("%f", math.Pow(2, conf.Game.Entity.MaxScale)),
				Cy:       "0.0",
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cx' Error:Field validation for 'cx' failed on the 'lte' tag"},
		}, {
			// top over
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", -math.Pow(2, conf.Game.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.cy' Error:Field validation for 'cy' failed on the 'gte' tag"},
		}, {
			// bottom over
			in: gameMapRequest{
				Cx:       "0.0",
				Cy:       fmt.Sprintf("%f", math.Pow(2, conf.Game.Entity.MaxScale)),
				Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
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
		if rawResult := v.Struct(c.in); rawResult == nil {
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

func TestGameMap(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cases := []struct {
			in   gameMapRequest
			want bool
		}{
			{
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MaxScale),
					Delegate: "0.0",
				},
				want: true,
			},
		}
		for _, c := range cases {
			w, _, r := prepare()
			r.GET("/gamemap", GameMap)
			assertOkResponse(t, paramAssertOk{
				Method: "GET",
				Path:   fmt.Sprintf("/gamemap?cx=%s&cy=%s&scale=%s&delegate=%s", c.in.Cx, c.in.Cy, c.in.Scale, c.in.Delegate),
				R:      r,
				W:      w,
				Jwt:    "",
				In:     c.in,
				Assert: func(got map[string]interface{}) {
					if len(got) == 0 {
						t.Errorf("/game got %v, want not empty", got)
					}
				},
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		cases := []struct {
			in   gameMapRequest
			want []string
		}{
			{
				// too small Scale
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale-0.0001),
					Delegate: "0.0",
				},
				want: []string{"scale must be gte 0.000000"},
			},
			{
				// too large Scale
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MaxScale+0.0001),
					Delegate: "0.0",
				},
				want: []string{"scale must be lte 16.000000"},
			},
			{
				// too small Delegate
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
					Delegate: "-0.0001",
				},
				want: []string{"delegate must be gte 0"},
			},
			{
				// too large Delegate
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
					Delegate: fmt.Sprintf("%f", conf.Game.Entity.MinScale+1),
				},
				want: []string{"delegate must be lte 0.000000"},
			},
			{
				// left over
				in: gameMapRequest{
					Cx:       fmt.Sprintf("%f", -math.Pow(2, conf.Game.Entity.MaxScale)),
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cx must be gte -32767.500000"},
			},
			{
				// right over
				in: gameMapRequest{
					Cx:       fmt.Sprintf("%f", math.Pow(2, conf.Game.Entity.MaxScale)),
					Cy:       "0.0",
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cx must be lte 32767.500000"},
			},
			{
				// top over
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       fmt.Sprintf("%f", -math.Pow(2, conf.Game.Entity.MaxScale)),
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
					Delegate: "0.0",
				},
				want: []string{"cy must be gte -32767.500000"},
			},
			{
				// bottom over
				in: gameMapRequest{
					Cx:       "0.0",
					Cy:       fmt.Sprintf("%f", math.Pow(2, conf.Game.Entity.MaxScale)),
					Scale:    fmt.Sprintf("%f", conf.Game.Entity.MinScale),
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
			w, _, r := prepare()
			r.GET("/gamemap", GameMap)
			url := fmt.Sprintf("/gamemap?cx=%s&cy=%s&scale=%s&delegate=%s", c.in.Cx, c.in.Cy, c.in.Scale, c.in.Delegate)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assertErrorResponse(url, t, w, c.want)
		}
	})
}
