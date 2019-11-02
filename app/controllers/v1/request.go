package v1

import (
	"fmt"
	"reflect"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// OwnerRequest represents requirement field of user action.
type OwnerRequest struct {
	O *entities.Player
}

// Parse validates and insert value from response.
func (o *OwnerRequest) Parse(token string) []string {
	errs := []string{}
	// [TODO] delete it
	//if o.O = services.FindOwner(token); o.O == nil {
	errs = append(errs, "user not found")
	//}
	return errs
}

// ScaleRequest represents requirement field of user action with scaling.
type ScaleRequest struct {
	OwnerRequest
	Scale float64
}

// Parse validates and insert value from response.
func (s *ScaleRequest) Parse(token string, v map[string]interface{}) []string {
	errs := s.OwnerRequest.Parse(token)
	if scale, ok := v["scale"].(float64); !ok {
		errs = append(errs, fmt.Sprintf("parse scale failed: %v", v["scale"]))
	} else {
		s.Scale = scale
	}
	return errs
}

// PointRequest represents requirement field of user action pointing somewhere.
type PointRequest struct {
	ScaleRequest
	X float64
	Y float64
}

// Parse validates and insert value from response.
func (p *PointRequest) Parse(token string, v map[string]interface{}) []string {
	errs := p.ScaleRequest.Parse(token, v)
	if x, ok := v["x"].(float64); !ok {
		errs = append(errs, fmt.Sprintf("parse x failed: %v", v["x"]))
	} else {
		p.X = x
	}
	if y, ok := v["y"].(float64); !ok {
		errs = append(errs, fmt.Sprintf("parse y failed: %v", v["y"]))
	} else {
		p.Y = y
	}
	return errs
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
