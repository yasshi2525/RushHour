package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/services"
)

// gameMapRequest represents requirement to view game map
type gameMapRequest struct {
	// X is center x scalized coordinate
	X string `form:"x" json:"x" validate:"required,numeric"`
	// Y is center y scalized coordinate
	Y string `form:"y" json:"y" validate:"required,numeric"`
	// Scale is 2^Scale coordinate maps size
	Scale string `form:"scale" json:"scale" validate:"required,numeric"`
	// Delegate is 2^Delegate grid of map
	Delegate string `form:"delegate" json:"delegate" validate:"required,numeric"`
}

// export converts string to float64
func (v *gameMapRequest) export() (int, int, int, int) {
	x, _ := strconv.ParseInt(v.X, 10, 64)
	y, _ := strconv.ParseInt(v.Y, 10, 64)
	sc, _ := strconv.ParseInt(v.Scale, 10, 64)
	dlg, _ := strconv.ParseInt(v.Delegate, 10, 64)
	return int(x), int(y), int(sc), int(dlg)
}

// validGameMapRequest validates that GameMapRequest contains game whole map
func validGameMapRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(gameMapRequest)
	x, y, sc, dlg := v.export()

	minSc := conf.Game.Entity.MinScale
	maxSc := conf.Game.Entity.MaxScale

	// validate scale
	if sc < minSc {
		sl.ReportError(v.Scale, "scale", "Scale", "gte", fmt.Sprintf("%d", minSc))
		return
	}
	if sc > maxSc {
		sl.ReportError(v.Scale, "scale", "Scale", "lte", fmt.Sprintf("%d", maxSc))
		return
	}

	// validate delegate
	if dlg < 0 {
		sl.ReportError(v.Delegate, "delegate", "Delegate", "gte", "0")
	}

	if sc-dlg < minSc {
		sl.ReportError(v.Delegate, "delegate", "Delegate", "lte", fmt.Sprintf("%d", sc-minSc))
	}

	length := 1 << sc
	border := 1 << maxSc

	// left over
	if x < 0 {
		sl.ReportError(v.X, "x", "X", "gte", fmt.Sprintf("%d", 0))
	}
	// right over
	if x+length > border {
		sl.ReportError(v.X, "x", "X", "lte", fmt.Sprintf("%d", border-length))
	}
	// top over
	if y < 0 {
		sl.ReportError(v.Y, "y", "Y", "gte", fmt.Sprintf("%d", 0))
	}
	// bottom over
	if y+length > border {
		sl.ReportError(v.Y, "y", "Y", "lte", fmt.Sprintf("%d", border-length))
	}
}

// GameMap returns all data of gamemap
// @Description entities are delegate object
// @Tags entities.DelegateMap
// @Summary get all entities in specified area
// @Accept json
// @Produce json
// @Param cx query number true "x coordinate"
// @Param cy query number true "y coordinate"
// @Param scale query number true "width,height(100%)=2^scale"
// @Param delegate query number true "width,height(grid)=2^delegate"
// @Success 200 {object} entities.DelegateMap "map centered (x,y) with grid in (width,height)"
// @Failure 400 {array} string "reasons of error when cx, cy and scale are out of area"
// @Router /gamemap [get]
func GameMap(c *gin.Context) {
	params := gameMapRequest{}
	if err := c.ShouldBind(&params); err != nil {
		c.Set(keyErr, err)
	} else {
		c.Set(keyOk, services.ViewDelegateMap(params.export()))
	}
}
