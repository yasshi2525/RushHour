package controllers

import (
	"math/rand"
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
	return c.RenderJSON(
		genResponse(
			true,
			services.ViewMap(
				rand.Float64()*100,
				rand.Float64()*100,
				rand.Float64()*7)))
}

// Diff returns only diff
func (c APIv1Game) Diff() revel.Result {
	return c.RenderJSON(
		genResponse(
			true,
			services.ViewMap(
				rand.Float64()*100,
				rand.Float64()*100,
				rand.Float64()*7,
				time.Now().Add(time.Duration(-1)*time.Minute))))
}

func genResponse(status bool, results interface{}) interface{} {
	return &struct {
		Status  bool        `json:"status"`
		Results interface{} `json:"results"`
	}{
		Status:  true,
		Results: results,
	}
}
