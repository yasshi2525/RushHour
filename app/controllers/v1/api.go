package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services"
)

// APIv1Game is controller for REST API
type APIv1Game struct {
	*revel.Controller
}

// Index returns gamemap
func (c APIv1Game) Index() revel.Result {
	var params = make(map[string]float64)

	for _, p := range []string{"cx", "cy", "scale"} {
		if v, err := strconv.ParseFloat(c.Params.Get(p), 64); err == nil {
			params[p] = v
		}
	}

	return c.RenderJSON(
		genResponse(true, services.ViewDelegateMap(params["cx"], params["cy"], params["scale"])))
}

// Diff returns only diff
func (c APIv1Game) Diff() revel.Result {
	return c.RenderJSON(
		genResponse(
			true,
			services.ViewMap(500, 500, 10, time.Now().Add(time.Duration(-1)*time.Minute))))
}

func (c APIv1Game) Departure() revel.Result {
	x, err := strconv.ParseFloat(c.Params.Form.Get("x"), 64)
	if err != nil {
		return c.RenderJSON(genResponse(false, err.Error()))
	}
	y, err := strconv.ParseFloat(c.Params.Form.Get("y"), 64)
	if err != nil {
		return c.RenderJSON(genResponse(false, err.Error()))
	}
	oid, err := strconv.ParseUint(c.Params.Form.Get("oid"), 10, 64)
	if err != nil {
		return c.RenderJSON(genResponse(false, err.Error()))
	}
	if o, ok := services.Model.Players[uint(oid)]; !ok {
		return c.RenderJSON(genResponse(false, fmt.Sprintf("%d not exists", oid)))
	} else {
		rn, err := services.CreateRailNode(o, x, y)
		if err != nil {
			return c.RenderJSON(genResponse(false, err.Error()))
		}
		return c.RenderJSON(genResponse(true, &struct {
			OwnerID uint    `json:"oid"`
			ID      uint    `json:"id"`
			X       float64 `json:"x"`
			Y       float64 `json:"y"`
		}{
			OwnerID: uint(oid),
			ID:      rn.ID,
			X:       rn.X,
			Y:       rn.Y,
		}))
	}
}

func genResponse(status bool, results interface{}) interface{} {
	return &struct {
		Status    bool        `json:"status"`
		Timestamp int64       `json:"timestamp"`
		Results   interface{} `json:"results"`
	}{
		Status:    true,
		Timestamp: time.Now().Unix(),
		Results:   results,
	}
}
