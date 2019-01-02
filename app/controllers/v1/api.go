package controllers

import (
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
	start := time.Now()
	defer services.WarnLongExec(start, 1, "JSON生成", true)

	r := struct {
		State   bool
		Results interface{}
	}{
		State:   true,
		Results: services.ViewMap(),
	}

	return c.RenderJSON(r)
}
