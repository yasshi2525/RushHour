package v1

import (
	"fmt"
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
				X:        "0.0",
				Y:        "0.0",
				Scale:    fmt.Sprintf("%d.0", conf.Game.Entity.MinScale),
				Delegate: "0.0",
			},
			want: nil,
		}, {
			// too small Scale
			in: gameMapRequest{
				X:        "0",
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale-1),
				Delegate: "0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'gte' tag"},
		}, {
			// too large Scale
			in: gameMapRequest{
				X:        "0",
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MaxScale+1),
				Delegate: "0.0",
			},
			want: []string{"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'lte' tag"},
		}, {
			// too small Delegate
			in: gameMapRequest{
				X:        "0",
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: "-1",
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'gte' tag"},
		}, {
			// too large Delegate
			in: gameMapRequest{
				X:        "0",
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: fmt.Sprintf("%d", conf.Game.Entity.MinScale+1),
			},
			want: []string{"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'lte' tag"},
		}, {
			// left over
			in: gameMapRequest{
				X:        "-1",
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: "0",
			},
			want: []string{"Key: 'gameMapRequest.x' Error:Field validation for 'x' failed on the 'gte' tag"},
		}, {
			// right over
			in: gameMapRequest{
				X:        fmt.Sprintf("%d", 0x1<<(conf.Game.Entity.MaxScale-conf.Game.Entity.MinScale)),
				Y:        "0",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: "0",
			},
			want: []string{"Key: 'gameMapRequest.x' Error:Field validation for 'x' failed on the 'lte' tag"},
		}, {
			// top over
			in: gameMapRequest{
				X:        "0",
				Y:        "-1",
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: "0",
			},
			want: []string{"Key: 'gameMapRequest.y' Error:Field validation for 'y' failed on the 'gte' tag"},
		}, {
			// bottom over
			in: gameMapRequest{
				X:        "0",
				Y:        fmt.Sprintf("%d", 0x1<<(conf.Game.Entity.MaxScale-conf.Game.Entity.MinScale)),
				Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
				Delegate: "0",
			},
			want: []string{"Key: 'gameMapRequest.y' Error:Field validation for 'y' failed on the 'lte' tag"},
		}, {
			// invalid format
			in: gameMapRequest{
				X:        "invalid",
				Y:        "invalid",
				Scale:    "invalid",
				Delegate: "invalid",
			},
			want: []string{
				"Key: 'gameMapRequest.x' Error:Field validation for 'x' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.y' Error:Field validation for 'y' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.scale' Error:Field validation for 'scale' failed on the 'numeric' tag",
				"Key: 'gameMapRequest.delegate' Error:Field validation for 'delegate' failed on the 'numeric' tag",
			},
		}, {
			// empty
			in: gameMapRequest{},
			want: []string{
				"Key: 'gameMapRequest.x' Error:Field validation for 'x' failed on the 'required' tag",
				"Key: 'gameMapRequest.y' Error:Field validation for 'y' failed on the 'required' tag",
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
				for i := 0; i < len(res); i++ {
					got := res[i]
					t.Errorf("validgameMapRequest(%v)[%d] == %s", c.in, i, got)
				}
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
					X:        "0.0",
					Y:        "0.0",
					Scale:    fmt.Sprintf("%d.0", conf.Game.Entity.MaxScale),
					Delegate: "0.0",
				},
				want: true,
			},
		}
		for _, c := range cases {
			w, _, r := prepare(ModelHandler())
			r.GET("/gamemap", GameMap)
			assertOkResponse(t, paramAssertOk{
				Method: "GET",
				Path:   fmt.Sprintf("/gamemap?x=%s&y=%s&scale=%s&delegate=%s", c.in.X, c.in.Y, c.in.Scale, c.in.Delegate),
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
					X:        "0",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale-1),
					Delegate: "0.0",
				},
				want: []string{"scale must be gte 0"},
			},
			{
				// too large Scale
				in: gameMapRequest{
					X:        "0",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MaxScale+1),
					Delegate: "0",
				},
				want: []string{"scale must be lte 16"},
			},
			{
				// too small Delegate
				in: gameMapRequest{
					X:        "0",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: "-1",
				},
				want: []string{"delegate must be gte 0"},
			},
			{
				// too large Delegate
				in: gameMapRequest{
					X:        "0",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: fmt.Sprintf("%d", conf.Game.Entity.MinScale+1),
				},
				want: []string{"delegate must be lte 0"},
			},
			{
				// left over
				in: gameMapRequest{
					X:        "-1",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: "0",
				},
				want: []string{"x must be gte 0"},
			},
			{
				// right over
				in: gameMapRequest{
					X:        "65536",
					Y:        "0",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: "0",
				},
				want: []string{"x must be lte 65535"},
			},
			{
				// top over
				in: gameMapRequest{
					X:        "0",
					Y:        "-1",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: "0",
				},
				want: []string{"y must be gte 0"},
			},
			{
				// bottom over
				in: gameMapRequest{
					X:        "0",
					Y:        "65536",
					Scale:    fmt.Sprintf("%d", conf.Game.Entity.MinScale),
					Delegate: "0",
				},
				want: []string{"y must be lte 65535"},
			},
			{
				// invalid format
				in: gameMapRequest{
					X:        "invalid",
					Y:        "invalid",
					Scale:    "invalid",
					Delegate: "invalid",
				},
				want: []string{
					"x must be numeric",
					"y must be numeric",
					"scale must be numeric",
					"delegate must be numeric",
				},
			},
			{
				// empty
				in: gameMapRequest{},
				want: []string{
					"x must be required",
					"y must be required",
					"scale must be required",
					"delegate must be required",
				},
			},
		}
		for _, c := range cases {
			w, _, r := prepare(ModelHandler())
			r.GET("/gamemap", GameMap)
			url := fmt.Sprintf("/gamemap?x=%s&y=%s&scale=%s&delegate=%s", c.in.X, c.in.Y, c.in.Scale, c.in.Delegate)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assertErrorResponse(url, t, w, c.want)
		}
	})
}
