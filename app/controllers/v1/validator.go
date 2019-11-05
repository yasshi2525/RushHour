package v1

// custom validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
)

// DefaultValidator enables custom validation
type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

// ValidateStruct is implementation of binding.StructValidator
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

// Engine is implementation of binding.StructValidator
func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = initValidate()
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func initValidate() *validator.Validate {
	v := validator.New()
	// BUG: err.Field() refered `json` field, but refer struct field
	// https://github.com/go-playground/validator/issues/337
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterStructValidation(validGameMapRequest, gameMapRequest{})
	v.RegisterStructValidation(validRegisterRequest, registerRequest{})
	v.RegisterStructValidation(validScaleRequest, scaleRequest{})
	v.RegisterStructValidation(validPointRequest, pointRequest{})
	return v
}
