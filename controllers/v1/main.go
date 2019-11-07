package v1

import (
	"fmt"
	"math"
	"reflect"

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
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
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
	Scale float64 `form:"scale" json:"scale" validate:"required"`
}

func validScaleRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(scaleRequest)

	minSc := conf.Game.Entity.MinScale
	maxSc := conf.Game.Entity.MaxScale

	// validate scale
	if v.Scale < minSc {
		sl.ReportError(v.Scale, "scale", "Scale", "gte", fmt.Sprintf("%f", minSc))
		return
	}
	if v.Scale > maxSc {
		sl.ReportError(v.Scale, "scale", "Scale", "lte", fmt.Sprintf("%f", maxSc))
		return
	}
}

type pointRequest struct {
	X     float64 `form:"x" json:"x" validate:"required"`
	Y     float64 `form:"y" json:"y" validate:"required"`
	Scale float64 `form:"scale" json:"scale" validate:"required"`
}

func validPointRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(pointRequest)

	minSc := conf.Game.Entity.MinScale
	maxSc := conf.Game.Entity.MaxScale

	// validate scale
	if v.Scale < minSc {
		sl.ReportError(v.Scale, "scale", "Scale", "gte", fmt.Sprintf("%f", minSc))
		return
	}
	if v.Scale > maxSc {
		sl.ReportError(v.Scale, "scale", "Scale", "lte", fmt.Sprintf("%f", maxSc))
		return
	}

	border := math.Pow(2, maxSc-1)

	// left over
	if v.X < -border {
		sl.ReportError(v.X, "cx", "Cx", "gte", fmt.Sprintf("%f", -border))
	}
	// right over
	if v.X > border {
		sl.ReportError(v.X, "cx", "Cx", "lte", fmt.Sprintf("%f", border))
	}
	// top over
	if v.Y < -border {
		sl.ReportError(v.Y, "cy", "Cy", "gte", fmt.Sprintf("%f", -border))
	}
	// bottom over
	if v.Y > border {
		sl.ReportError(v.Y, "cy", "Cy", "lte", fmt.Sprintf("%f", border))
	}
}

func validateEntity(res entities.ModelType, id uint) (entities.Entity, error) {
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
