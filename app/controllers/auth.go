package controllers

import (
	"github.com/gomodule/oauth1/oauth"
	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// Auth is revel controller handling authentication
type Auth struct {
	*revel.Controller
}

// Twitter redirects twitter sign in page
func (a Auth) Twitter() revel.Result {
	tmpCred, url := auth.GetTwitterAuthURL()
	a.Session.Set("tmpToken", tmpCred.Token)
	a.Session.Set("tmpSecret", tmpCred.Secret)
	return a.Redirect(url)
}

// TwitterCallback registers user info
func (a Auth) TwitterCallback() revel.Result {
	secret := a.Params.Get("oauth_verifier")

	if tmpToken, err := a.Session.Get("tmpToken"); err != nil {
		return a.RenderHTML(err.Error())
	} else {
		if tmpSecret, err := a.Session.Get("tmpSecret"); err != nil {
			return a.RenderHTML(err.Error())
		} else {
			if info, err := auth.GetTwitterUserInfo(&oauth.Credentials{
				Token:  tmpToken.(string),
				Secret: tmpSecret.(string),
			}, secret); err != nil {
				return a.RenderHTML(err.Error())
			} else {
				if my, err := services.OAuthSignIn(entities.Twitter, info); err != nil {
					return a.RenderHTML(err.Error())
				} else {
					a.Session.Set("token", my.Token)
					return a.Redirect("/")
				}
			}
		}
	}
}

// Google redirects google sign in page
func (a Auth) Google() revel.Result {
	return a.Redirect(auth.GetGoogleAuthURL())
}

// GoogleCallback registers user info
func (a Auth) GoogleCallback() revel.Result {
	state := a.Params.Get("state")
	code := a.Params.Get("code")

	if info, err := auth.GetGoogleUserInfo(state, code); err != nil {
		return a.RenderHTML(err.Error())
	} else {
		if my, err := services.OAuthSignIn(entities.Google, info); err != nil {
			return a.RenderHTML(err.Error())
		} else {
			a.Session.Set("token", my.Token)
		}
		return a.Redirect("/")
	}
}

// GitHub redirects google sign in page
func (a Auth) GitHub() revel.Result {
	return a.Redirect(auth.GetGitHubAuthURL())
}

// GitHubCallback registers user info
func (a Auth) GitHubCallback() revel.Result {
	state := a.Params.Get("state")
	code := a.Params.Get("code")

	if info, err := auth.GetGitHubUserInfo(state, code); err != nil {
		return a.RenderHTML(err.Error())
	} else {
		if my, err := services.OAuthSignIn(entities.GitHub, info); err != nil {
			return a.RenderHTML(err.Error())
		} else {
			a.Session.Set("token", my.Token)
		}
		return a.Redirect("/")
	}
}
