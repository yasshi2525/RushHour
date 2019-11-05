package v1

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/app/services"
)

const getGameMapRequest string = "dive,keys,oneof=cx cy scale delegate,endkeys"

// gameMapRequest represents requirement to view game map
type gameMapRequest struct {
	// Cx is center x coordinate
	Cx string `form:"cx" json:"cx" validate:"required,numeric"`
	// Cy is center y coordinate
	Cy string `form:"cy" json:"cy" validate:"required,numeric"`
	// Scale is 2^Scale coordinate maps size
	Scale string `form:"scale" json:"scale" validate:"required,numeric"`
	// Delegate is 2^Delegate grid of map
	Delegate string `form:"delegate" json:"delegate" validate:"required,numeric"`
}

// export converts string to float64
func (v *gameMapRequest) export() (float64, float64, float64, float64) {
	cx, _ := strconv.ParseFloat(v.Cx, 64)
	cy, _ := strconv.ParseFloat(v.Cy, 64)
	sc, _ := strconv.ParseFloat(v.Scale, 64)
	dlg, _ := strconv.ParseFloat(v.Delegate, 64)
	return cx, cy, sc, dlg
}

// validGameMapRequest validates that GameMapRequest contains game whole map
func validGameMapRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(gameMapRequest)
	cx, cy, sc, dlg := v.export()

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

	// validate delegate
	if dlg < 0 {
		sl.ReportError(v.Delegate, "delegate", "Delegate", "gte", "0")
	}

	if sc-dlg < minSc {
		sl.ReportError(v.Delegate, "delegate", "Delegate", "lte", fmt.Sprintf("%f", sc-minSc))
	}

	radius := math.Pow(2, sc-1)
	border := math.Pow(2, maxSc-1)

	// left over
	if cx-radius < -border {
		sl.ReportError(v.Cx, "cx", "Cx", "gte", fmt.Sprintf("%f", radius-border))
	}
	// right over
	if cx+radius > border {
		sl.ReportError(v.Cx, "cx", "Cx", "lte", fmt.Sprintf("%f", border-radius))
	}
	// top over
	if cy-radius < -border {
		sl.ReportError(v.Cx, "cy", "Cy", "gte", fmt.Sprintf("%f", radius-border))
	}
	// bottom over
	if cy+radius > border {
		sl.ReportError(v.Cx, "cy", "Cy", "lte", fmt.Sprintf("%f", border-radius))
	}
}

// GameMap returns all data of gamemap
// @Description entities are delegate object
// @Tags entities.DelegateMap
// @Summary get all entities in specified area
// @Accept  json
// @Produce  json
// @Param cx query number true "x coordinate"
// @Param cy query number true "y coordinate"
// @Param scale query number true "width,height(100%)=2^scale"
// @Param delegate query number true "width,height(grid)=2^delegate"
// @Success 200 {object} entities.DelegateMap "map centered (x,y) with grid in (width,height)"
// @Failure 400 {array} string "reasons of error when cx, cy and scale are out of area"
// @Router /gamemap [get]
func GameMap(c *gin.Context) {
	params := gameMapRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else {
		c.Set(keyOk, services.ViewDelegateMap(params.export()))
	}
}
