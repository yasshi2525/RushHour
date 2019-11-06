package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
	"github.com/yasshi2525/RushHour/services"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	conf = &config.Config{}
	conf.Game.Entity.MinScale = 0
	conf.Game.Entity.MaxScale = 16
	conf.Game.Service.Procedure.Interval.D = 1 * time.Hour
	conf.Secret.Admin.UserName = "admin@example.com"
	conf.Secret.Admin.Password = "password_test"
	var err error
	if auther, err = auth.GetAuther(conf.Secret.Auth); err != nil {
		panic(fmt.Errorf("failed to create auther: %v", err))
	} else {
		services.Init(&services.ServiceConfig{
			AppConf: conf,
			Auther:  auther,
		})
		defer services.Terminate()
		services.CreateIfAdmin()
		services.Start()
		defer services.Stop()
		binding.Validator = new(DefaultValidator)
		m.Run()
	}
}

func prepare(handlers ...gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.Use(handlers...)
	return w, c, r
}

func assertValidation(name string, t *testing.T, v *validator.Validate, in interface{}, want []string) {
	t.Helper()
	if rawResult := v.Struct(in); rawResult == nil {
		if want != nil {
			t.Errorf("%s(%v) got nil, want %v", name, in, want)
		}
	} else {
		res := rawResult.(validator.ValidationErrors)
		if len(res) != len(want) {
			t.Errorf("%s(%v).err got %d errors, want %d errors", name, in, len(res), len(want))
		} else {
			for i := 0; i < len(res); i++ {
				got, want := res[i], want[i]
				if fmt.Sprintf("%s", got) != want {
					t.Errorf("%s(%v)[%d] got %s, want %s", name, in, i, got, want)
				}
			}
		}
	}
}

func assertErrorResponse(name string, t *testing.T, w *httptest.ResponseRecorder, want []string) {
	t.Helper()
	if w.Code != http.StatusBadRequest {
		t.Errorf("%s.code got %d, want %d (details = %s)", name, w.Code, http.StatusBadRequest, w.Body.String())
	} else {
		var msg errInfo
		if err := json.Unmarshal(w.Body.Bytes(), &msg); err != nil {
			t.Errorf("%s.body.err got %v, want nil", name, err)
		} else {
			got := msg.Err
			if len(got) != len(want) {
				t.Errorf("%s.body got %d errors, want %d errors", name, len(got), len(want))
			} else {
				for i := 0; i < len(got); i++ {
					if got, want := got[i], want[i]; got != want {
						t.Errorf("%s.body[%d] got %s, want %s", name, i, got, want)
					}
				}
			}
		}
	}
}

type paramAssertOk struct {
	Method string
	Path   string
	Jwt    string
	R      *gin.Engine
	W      *httptest.ResponseRecorder
	In     interface{}
	Assert func(got map[string]interface{})
}

func assertOkResponse(t *testing.T, args paramAssertOk) {
	t.Helper()
	str, _ := json.Marshal(args.In)
	req, _ := http.NewRequest(args.Method, args.Path, bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")
	if args.Jwt != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", args.Jwt))
	}
	args.R.ServeHTTP(args.W, req)
	if args.W.Code != http.StatusOK {
		t.Errorf("%s.code got %d, want %d (details = %s)", args.Path, args.W.Code, http.StatusOK, args.W.Body.String())
	} else {
		var msg map[string]interface{}
		if err := json.Unmarshal(args.W.Body.Bytes(), &msg); err != nil {
			t.Errorf("%s.body.err got %v, want nil", args.Path, err)
		} else {
			args.Assert(msg)
		}
	}
}

type paramAssertUnauthorized struct {
	Method string
	Path   string
	R      *gin.Engine
	W      *httptest.ResponseRecorder
	In     interface{}
}

func assertUnauthorized(t *testing.T, args paramAssertUnauthorized) {
	t.Helper()
	str, _ := json.Marshal(args.In)
	req, _ := http.NewRequest(args.Method, args.Path, bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")
	args.R.ServeHTTP(args.W, req)
	if args.W.Code != http.StatusUnauthorized {
		t.Errorf("%s.code got %d, want %d (details = %s)", args.Path, args.W.Code, http.StatusUnauthorized, args.W.Body.String())
	}
}

func TestInitController(t *testing.T) {
	bkConf := conf
	wantConf := &config.Config{
		Game: config.CnfGame{
			Entity: config.CnfEntity{
				MaxScale: 5,
			},
		},
	}
	wantAuther, _ := auth.GetAuther(wantConf.Secret.Auth)

	InitController(wantConf, wantAuther)
	if conf != wantConf {
		t.Errorf("InitController(%v, %v).conf got %v, want %v", wantConf, wantAuther, conf, wantConf)
	}
	if auther != wantAuther {
		t.Errorf("InitController(%v, %v).auther got %v, want %v", wantConf, wantAuther, auther, wantAuther)
	}
	conf = bkConf
}

func registerTestUser(t *testing.T, id string, password string) string {
	t.Helper()
	w, _, r := prepare()
	r.POST("/register", Register)
	str, _ := json.Marshal(&registerRequest{
		loginRequest: loginRequest{ID: id, Password: password},
		DisplayName:  id,
		Hue:          "0",
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if token, ok := res["jwt"]; ok {
		return token.(string)
	}
	t.Fatalf("register failed: %v", res)
	return ""
}
