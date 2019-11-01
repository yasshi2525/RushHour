package v1

import (
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

// Response is general response format when return code is 200
type Response struct {
	Status    bool        `json:"status"`
	Timestamp int64       `json:"timestamp"`
	Results   interface{} `json:"results"`
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
	validate.RegisterStructValidation(validGameMapRequest, GameMapRequest{})
}
