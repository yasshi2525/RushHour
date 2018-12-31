package controllers

import (
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services"
)

// Response represents general structure of Rest
type Response struct {
	State   bool
	Results interface{}
}

// ApiV1Game is controller for REST API
type ApiV1Game struct {
	*revel.Controller
}

// Index returns gamemap
func (c ApiV1Game) Index() revel.Result {

	r := Response{
		State:   true,
		Results: services.ViewMap(),
	}

	return c.RenderJSON(r)
}
