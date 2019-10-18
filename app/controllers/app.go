package controllers

import (
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// App is controller for index.html
type App struct {
	*revel.Controller
}

// Index renders index.html
func (c App) Index() revel.Result {
	c.ViewArgs["loggedin"] = false
	if token, err := c.Session.Get("token"); err == nil {
		if o := services.FindOwner(token.(string)); o != nil {
			c.ViewArgs["loggedin"] = true
			c.ViewArgs["oid"] = o.ID
			info := o.ExportInfo()
			c.ViewArgs["name"] = info.DisplayName
			c.ViewArgs["image"] = info.Image
			c.ViewArgs["hue"] = o.Hue
			c.ViewArgs["admin"] = o.Level == entities.Admin
		} else {
			c.Session.Del("token")
		}
	} else {
		c.Session.Del("token")
	}
	c.ViewArgs["maintenance"] = !services.IsInOperation()
	return c.Render()
}

// IndexPost redirect to index.html
func (c App) IndexPost() revel.Result {
	return c.Redirect("/")
}

// SignOut delete session attribute token.
func (c App) SignOut() revel.Result {
	if token, err := c.Session.Get("token"); err == nil {
		services.SignOut(token.(string))
		c.Session.Del("token")
	}
	return c.Redirect("/")
}
