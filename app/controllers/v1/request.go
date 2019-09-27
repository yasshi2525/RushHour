package v1

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// OwnerRequest represents requirement field of user action.
type OwnerRequest struct {
	O     *entities.Player
	Scale float64
}

// Parse validates and insert value from response.
func (o *OwnerRequest) Parse(token string, v url.Values) []string {
	errs := []string{}

	if o.O = services.FindOwner(token); o.O == nil {
		errs = append(errs, "user not found")
	}
	if scale, err := strconv.ParseFloat(v.Get("scale"), 64); err != nil {
		errs = append(errs, err.Error())
	} else {
		o.Scale = scale
	}
	return errs
}

// PointRequest represents requirement field of user action pointing somewhere.
type PointRequest struct {
	OwnerRequest
	X float64
	Y float64
}

// Parse validates and insert value from response.
func (p *PointRequest) Parse(token string, v url.Values) []string {
	errs := p.OwnerRequest.Parse(token, v)
	if x, err := strconv.ParseFloat(v.Get("x"), 64); err != nil {
		errs = append(errs, err.Error())
	} else {
		p.X = x
	}
	if y, err := strconv.ParseFloat(v.Get("y"), 64); err != nil {
		errs = append(errs, err.Error())
	} else {
		p.Y = y
	}
	return errs
}

func validateEntity(res entities.ModelType, idstr string) (entities.Entity, error) {
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, err
	}
	val := services.Model.Values[res].MapIndex(reflect.ValueOf(uint(id)))
	if !val.IsValid() {
		return nil, fmt.Errorf("%s[%d] doesn't exist", res.String(), id)
	}
	return val.Interface().(entities.Entity), nil
}
