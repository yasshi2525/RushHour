package controllers

import (
	"math/rand"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services"
)

// APIv1Game is controller for REST API
type APIv1Game struct {
	*revel.Controller
}

// Index returns gamemap
func (c APIv1Game) Index() revel.Result {
	r := struct {
		Status  bool        `json:"status"`
		Results interface{} `json:"results"`
	}{
		Status: true,
		Results: services.ViewMap(
			rand.Float64()*100,
			rand.Float64()*100,
			rand.Float64()*7),
	}

	return c.RenderJSON(r)
}
