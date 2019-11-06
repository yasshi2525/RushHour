package controllers

import (
	"github.com/gomodule/oauth1/oauth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/app/entities"
)

// Twitter redirects twitter sign in page
func Twitter(c *gin.Context) {
	if cred, url, err := auther.GetTwitterAuthURL(); err != nil {
		c.Set(keyErr, err)
	} else {
		session := sessions.Default(c)
		session.Set("tmpToken", cred.Token)
		session.Set("tmpSecret", cred.Secret)
		session.Save()
		c.Set(keyRedirect, url)
	}
}

// TwitterCallback registers user info
// @Param oauth_verifier query string true "access token"
// @Failure 503 {string} string "session expire"
// @Failure 503 {string} string "no OAuth token"
func TwitterCallback(c *gin.Context) {
	session := sessions.Default(c)
	secret := c.Query("oauth_verifier")

	if tmpToken := session.Get("tmpToken"); tmpToken == nil {
		c.Set(keyErr, "no session key tmpToken")
	} else {
		if tmpSecret := session.Get("tmpSecret"); tmpSecret == nil {
			c.Set(keyErr, "no session key tmpSecret")
		} else {
			if info, err := auther.GetTwitterOAuthInfo(&oauth.Credentials{
				Token:  tmpToken.(string),
				Secret: tmpSecret.(string),
			}, secret); err != nil {
				c.Set(keyErr, err)
			} else {
				c.Set(keyAuthType, entities.Twitter)
				c.Set(keyAuthInfo, info)
			}
		}
	}
}

// Google redirects google sign in page
func Google(c *gin.Context) {
	c.Set(keyRedirect, auther.GetGoogleAuthURL())
}

// GoogleCallback registers user info
// @Param state query string true "state"
// @Param code query string true "code"
// @Failure 503 {string} string "no OAuth token"
func GoogleCallback(c *gin.Context) {
	c.Set(keyAuthType, entities.Google)
	c.Set(keyAuthFunc, auther.GetGoogleOAuthInfo)
}

// GitHub redirects github sign in page
func GitHub(c *gin.Context) {
	c.Set(keyRedirect, auther.GetGitHubAuthURL())
}

// GitHubCallback registers user info
// @Description try register github info
// @Summary try register github info
// @Param state query string true "state"
// @Param code query string true "code"
// @Failure 503 {string} string "no OAuth token"
func GitHubCallback(c *gin.Context) {
	c.Set(keyAuthType, entities.GitHub)
	c.Set(keyAuthFunc, auther.GetGitHubOAuthInfo)
}
