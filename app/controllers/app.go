package controllers

import (
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services"
)

// App is controller for index.html
type App struct {
	*revel.Controller
}

// Index renders index.html
func (c App) Index() revel.Result {
	c.ViewArgs["inOperation"] = services.IsInOperation()
	return c.Render()
}

// IndexPost redirect to index.html
func (c App) IndexPost() revel.Result {
	return c.Redirect("/")
}
