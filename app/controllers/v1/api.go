package controllers

import (
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
