package v1

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/services"
)

// @title RushHour REST API
// @version 1.0
// @description RushHour REST API
// @license.name MIT License
// @host rushhourgame.net
// @BasePath /api/v1
// @schemes https

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
	Err []string `json:"err"`
}

var conf *config.Config
var auther *auth.Auther

// excludeList accepts admin request even under maintenance
var excludeList []string

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
func InitController(c *config.Config, a *auth.Auther) {
	conf = c
	auther = a
	excludeList = []string{
		"/api/v1/login",
		"/api/v1/admin",
	}
}
