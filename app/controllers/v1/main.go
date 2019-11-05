package v1

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// entry represents generic key-value pair
type entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type jwtInfo struct {
	Jwt string `json:"jwt"`
}

// user represents public attributes that everyone can view
type user struct {
	// ID is number
	ID uint `json:"id"`
	// Name is display name
	Name string `json:"name"`
	// Image is url of profile icon
	Image string `json:"image"`
	// Hue is rail line color (HSV model)
	Hue float64 `json:"hue"`
}

type errInfo struct {
	Err interface{} `json:"err"`
}

var conf config.Config
var auther *auth.Auther

// buildJwt returns JSON Web Token of player
func buildJwt(o *entities.Player) (*jwtInfo, error) {
	url := conf.Secret.Auth.BaseURL
	now := time.Now()
	exp := now.Add(time.Hour)
	uu := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":                        url,
		"sub":                        "AccessToken",
		"aud":                        url,
		"exp":                        exp.Unix(),
		"nbf":                        now.Unix(),
		"iat":                        now.Unix(),
		"jti":                        uu.String(),
		fmt.Sprintf("%s/id", url):    o.ID,
		fmt.Sprintf("%s/name", url):  auther.Decrypt(o.GetDisplayName()),
		fmt.Sprintf("%s/image", url): auther.Decrypt(o.GetImage()),
		fmt.Sprintf("%s/admin", url): o.Level == entities.Admin,
		fmt.Sprintf("%s/hue", url):   o.Hue,
	})

	jwt, err := token.SignedString([]byte(conf.Secret.Auth.Salt))
	if err != nil {
		return nil, err
	}
	return &jwtInfo{jwt}, nil
}

type scaleRequest struct {
	Scale string `form:"scale" json:"scale" validate:"required,numeric"`
}

func (v *scaleRequest) export() float64 {
	sc, _ := strconv.ParseFloat(v.Scale, 64)
	return sc
}

func validScaleRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(scaleRequest)
	sc := v.export()

	minSc := conf.Game.Entity.MinScale
	maxSc := conf.Game.Entity.MaxScale

	// validate scale
	if sc < minSc {
		sl.ReportError(v.Scale, "scale", "Scale", "gte", fmt.Sprintf("%f", minSc))
		return
	}
	if sc > maxSc {
		sl.ReportError(v.Scale, "scale", "Scale", "lte", fmt.Sprintf("%f", maxSc))
		return
	}
}

type pointRequest struct {
	X     string `form:"x" json:"x" validate:"required,numeric"`
	Y     string `form:"y" json:"y" validate:"required,numeric"`
	Scale string `form:"scale" json:"scale" validate:"required,numeric"`
}

func (v *pointRequest) export() (float64, float64, float64) {
	x, _ := strconv.ParseFloat(v.X, 64)
	y, _ := strconv.ParseFloat(v.Y, 64)
	sc, _ := strconv.ParseFloat(v.Scale, 64)
	return x, y, sc
}

func validPointRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(pointRequest)
	x, y, sc := v.export()

	minSc := conf.Game.Entity.MinScale
	maxSc := conf.Game.Entity.MaxScale

	// validate scale
	if sc < minSc {
		sl.ReportError(v.Scale, "scale", "Scale", "gte", fmt.Sprintf("%f", minSc))
		return
	}
	if sc > maxSc {
		sl.ReportError(v.Scale, "scale", "Scale", "lte", fmt.Sprintf("%f", maxSc))
		return
	}

	border := math.Pow(2, maxSc-1)

	// left over
	if x < -border {
		sl.ReportError(v.X, "cx", "Cx", "gte", fmt.Sprintf("%f", -border))
	}
	// right over
	if x > border {
		sl.ReportError(v.X, "cx", "Cx", "lte", fmt.Sprintf("%f", border))
	}
	// top over
	if y < -border {
		sl.ReportError(v.Y, "cy", "Cy", "gte", fmt.Sprintf("%f", -border))
	}
	// bottom over
	if y > border {
		sl.ReportError(v.Y, "cy", "Cy", "lte", fmt.Sprintf("%f", border))
	}
}

func validateEntity(res entities.ModelType, raw interface{}) (entities.Entity, error) {
	idnum, ok := raw.(float64)
	if !ok {
		return nil, fmt.Errorf("%s[%v] doesn't exist", res.String(), raw)
	}
	id := uint(idnum)
	val := services.Model.Values[res].MapIndex(reflect.ValueOf(id))
	if !val.IsValid() {
		return nil, fmt.Errorf("%s[%d] doesn't exist", res.String(), id)
	}
	return val.Interface().(entities.Entity), nil
}

// InitController loads config
func InitController(c config.Config, a *auth.Auther) {
	conf = c
	auther = a
}
