package v1

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

func TestValidRegisterRequest(t *testing.T) {
	Init()

	cases := []struct {
		in   registerRequest
		want []string
	}{
		{
			in: registerRequest{
				loginRequest: loginRequest{
					ID:       "test@example.com",
					Password: "password",
				},
				DisplayName: "Test",
				Hue:         "0",
			},
			want: nil,
		}, {
			// too small hue
			in: registerRequest{
				loginRequest: loginRequest{
					ID:       "test@example.com",
					Password: "password",
				},
				DisplayName: "",
				Hue:         "0",
			},
			want: []string{"Key: 'registerRequest.hue' Error:Field validation for 'hue' failed on the 'gte' tag"},
		}, {
			// too large hue
			in: registerRequest{
				loginRequest: loginRequest{
					ID:       "test@example.com",
					Password: "password",
				},
				DisplayName: "",
				Hue:         "360",
			},
			want: []string{"Key: 'registerRequest.hue' Error:Field validation for 'hue' failed on the 'lt' tag"},
		},
	}
	for _, c := range cases {
		if rawResult := validate.Struct(c.in); rawResult == nil {
			if c.want != nil {
				t.Errorf("registerRequest(%v) == nil, want %v", c.in, c.want)
			}
		} else {
			res := rawResult.(validator.ValidationErrors)

			if len(res) != len(c.want) {
				t.Errorf("registerRequest(%v) == %d errors, want %d errors", c.in, len(res), len(c.want))
			} else {
				for i := 0; i < len(res); i++ {
					got, want := res[i], c.want[i]
					if fmt.Sprintf("%s", got) != want {
						t.Errorf("registerRequest(%v)[%d] == %s, want %s", c.in, i, got, want)
					}
				}
			}
		}
	}
}
