package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"

	"github.com/yasshi2525/RushHour/services"
)

func TestLogin(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cases := []struct {
			in      loginRequest
			wantNot string
		}{
			{
				in: loginRequest{
					ID:       conf.Secret.Admin.UserName,
					Password: conf.Secret.Admin.Password,
				},
				wantNot: "",
			},
		}

		for _, c := range cases {
			w, _, r := prepare(ModelHandler())
			r.POST("/login", Login)
			assertOkResponse(t, paramAssertOk{
				Method: "POST",
				Path:   "/login",
				Jwt:    "",
				R:      r,
				W:      w,
				In:     c.in,
				Assert: func(got map[string]interface{}) {
					if got["jwt"] == c.wantNot {
						t.Errorf("/login.jwt got %s, not want %s", got, c.wantNot)
					} else {
						_, err := jwt.Parse(got["jwt"].(string), func(token *jwt.Token) (interface{}, error) {
							if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
								return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
							}
							return []byte(conf.Secret.Auth.Salt), nil
						})
						if err != nil {
							t.Errorf("/login.jwt.err got %v, want nil", err)
						}
					}
				},
			})
		}
	})
	t.Run("error", func(t *testing.T) {
		cases := []struct {
			in   loginRequest
			want []string
		}{
			{
				// invalid mail address
				in: loginRequest{
					ID:       "nobody",
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
			w, _, r := prepare(ModelHandler())
			r.POST("/login", Login)
			str, _ := json.Marshal(c.in)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(str))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assertErrorResponse("/login", t, w, c.want)
		}
	})
}

func TestValidRegisterRequest(t *testing.T) {
	v := initValidate()

	cases := []struct {
		in   registerRequest
		want []string
	}{
		{
			in: registerRequest{
				DisplayName: "Test",
				Hue:         0,
			},
			want: nil,
		}, {
			// too small hue
			in: registerRequest{
				DisplayName: "",
				Hue:         -1,
			},
			want: []string{"Key: 'registerRequest.hue' Error:Field validation for 'hue' failed on the 'gte' tag"},
		}, {
			// too large hue
			in: registerRequest{
				DisplayName: "",
				Hue:         360,
			},
			want: []string{"Key: 'registerRequest.hue' Error:Field validation for 'hue' failed on the 'lt' tag"},
		},
	}
	for _, c := range cases {
		assertValidation("registerRequest", t, v, c.in, c.want)
	}
}

func TestRegister(t *testing.T) {

	type actualRegisterReq struct {
		ID          string `json:"id"`
		Password    string `json:"password"`
		DisplayName string `json:"name"`
		Hue         int    `json:"hue"`
	}

	t.Run("ok", func(t *testing.T) {
		cases := []struct {
			in      actualRegisterReq
			wantNot string
		}{
			{
				in: actualRegisterReq{
					ID:          "ok@example.com",
					Password:    "password",
					DisplayName: "ok",
					Hue:         0,
				},
				wantNot: "",
			},
		}

		for _, c := range cases {
			w, _, r := prepare(ModelHandler())
			r.POST("/register", Register)
			assertOkResponse(t, paramAssertOk{
				Method: "POST",
				Path:   "/register",
				R:      r,
				W:      w,
				In:     c.in,
				Assert: func(got map[string]interface{}) {
					if got["jwt"] == c.wantNot {
						t.Errorf("/login.jwt got %s, not want %s", got["jwt"], c.wantNot)
					} else {
						_, err := jwt.Parse(got["jwt"].(string), func(token *jwt.Token) (interface{}, error) {
							if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
								return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
							}
							return []byte(conf.Secret.Auth.Salt), nil
						})
						if err != nil {
							t.Errorf("/login.jwt.err got %v, want nil", err)
						}
					}
				},
			})
		}
	})
	t.Run("error", func(t *testing.T) {
		cases := []struct {
			in   actualRegisterReq
			want []string
		}{
			{
				// too small hue
				in: actualRegisterReq{
					ID:          "test@example.com",
					Password:    "password",
					DisplayName: "",
					Hue:         -1,
				},
				want: []string{"hue must be gte 0"},
			},
			{
				// too large hue
				in: actualRegisterReq{
					ID:          "test@example.com",
					Password:    "password",
					DisplayName: "",
					Hue:         360,
				},
				want: []string{"hue must be lt 360"},
			},
		}
		for _, c := range cases {
			w, _, r := prepare(ModelHandler())
			r.POST("/register", Register)
			str, _ := json.Marshal(c.in)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(str))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assertErrorResponse("/register", t, w, c.want)
		}
	})
}

func TestSettings(t *testing.T) {
	token := registerTestUser(t, "setting@example.com", "password")

	t.Run("ok", func(t *testing.T) {
		cases := []struct {
			in   string
			want map[string]interface{}
		}{
			{
				in: token,
				want: map[string]interface{}{
					"email":        "setting@example.com",
					"custom_name":  "setting@example.com",
					"custom_image": "",
					"auth_type":    "RushHour",
				},
			},
		}
		for _, c := range cases {
			w, _, r := prepare(JWTHandler(), ModelHandler())
			r.GET("/settings", Settings)
			assertOkResponse(t, paramAssertOk{
				Method: "GET",
				Path:   "/settings",
				Jwt:    token,
				R:      r,
				W:      w,
				In:     "",
				Assert: func(got map[string]interface{}) {
					for key, want := range c.want {
						if got[key] != want {
							t.Errorf("/settings[%s] got %v, want %v", key, got[key], want)
						}
					}
				},
			})
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		w, _, r := prepare(JWTHandler(), ModelHandler())
		r.GET("/settings", Settings)
		assertUnauthorized(t, paramAssertUnauthorized{
			Method: "GET",
			Path:   "/settings",
			R:      r,
			W:      w,
			In:     "",
		})
	})
}

func fetchSettings(t *testing.T, token string) *services.AccountSettings {
	t.Helper()
	w, _, r := prepare(JWTHandler(), ModelHandler())
	r.GET("/settings", Settings)
	req, _ := http.NewRequest("GET", "/settings", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	as := services.AccountSettings{}
	if err := json.Unmarshal(w.Body.Bytes(), &as); err != nil {
		t.Fatalf("failed to get settings: %v", err)
		return nil
	}
	return &as
}

func TestChangeSettings(t *testing.T) {
	token := registerTestUser(t, "change_setting@example.com", "password")
	t.Run("ok", func(t *testing.T) {
		cases := []struct {
			in   entry
			want interface{}
		}{
			{
				in:   entry{Key: "custom_name", Value: "changed"},
				want: "changed",
			}, {
				in:   entry{Key: "use_cname", Value: true},
				want: true,
			},
		}
		for _, c := range cases {
			w, _, r := prepare(JWTHandler(), ModelHandler())
			r.POST("/settings/:resname", ChangeSettings)
			assertOkResponse(t, paramAssertOk{
				Method: "POST",
				Path:   fmt.Sprintf("/settings/%s", c.in.Key),
				Jwt:    token,
				R:      r,
				W:      w,
				In: struct {
					Value interface{} `json:"value"`
				}{c.in.Value},
				Assert: func(got map[string]interface{}) {
					if got["key"] != c.in.Key {
						t.Errorf("/settings/%s[key] got %v, want %v", c.in.Key, got["key"], c.in.Key)
					}
					if got["value"] != c.want {
						t.Errorf("/settings/%s got %v, want %v", c.in.Key, got["value"], c.want)
					}
				},
			})
		}
	})
}
