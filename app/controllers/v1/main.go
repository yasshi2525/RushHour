package v1

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/revel/revel"
	"gopkg.in/go-playground/validator.v9"
)

// @title RushHour REST API
// @version 1.0
// @description RushHour REST API
// @license.name MIT License
// @host rushhourgame.net
// @BasePath /api/v1
// @schemes https

// API is controller of REST server
type API struct {
	*revel.Controller
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

// mapToStruct converts request parameter map to struct filterd field key.
// key mapped `json:"<*>"` tag
// outPtr must be pointer. It returns outPtr.
func mapToStruct(origin url.Values, outPtr interface{}) interface{} {
	obj := reflect.ValueOf(outPtr).Elem()
	t := obj.Type()
	for i := 0; i < t.NumField(); i++ {
		obj.Field(i).Set(reflect.ValueOf(origin.Get(t.Field(i).Tag.Get("json"))))
	}
	return outPtr
}

func buildErrorMessages(errs validator.ValidationErrors) []string {
	msgs := []string{}
	for _, err := range errs {
		if err.Param() == "" {
			msgs = append(msgs, fmt.Sprintf("%s must be %s", err.Field(), err.Tag()))
		} else {
			msgs = append(msgs, fmt.Sprintf("%s must be %s %s", err.Field(), err.Tag(), err.Param()))
		}

	}
	return msgs
}

var validate *validator.Validate

// Init must be called in StartUp phase
func Init() {
	validate = validator.New()
	// BUG: err.Field() refered `json` field, but refer struct field
	// https://github.com/go-playground/validator/issues/337
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	validate.RegisterStructValidation(validGameMapRequest, gameMapRequest{})
}
