package v1

import (
	"net/http"

	"github.com/gomodule/oauth1/oauth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/services"
)

// Twitter redirects twitter sign in page
func Twitter(c *gin.Context) {
	if cred, url, err := auther.GetTwitterAuthURL(); err != nil {
		c.Set(keyErr, err.Error())
	} else {
		session := sessions.Default(c)
		session.Set("tmpToken", cred.Token)
		session.Set("tmpSecret", cred.Secret)
		session.Save()
		c.Redirect(http.StatusFound, url)
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
				c.Set(keyErr, err.Error())
			} else {
				if _, err := services.OAuthSignIn(entities.Twitter, info); err != nil {
					c.Set(keyErr, err.Error())
				} else {
					c.Redirect(http.StatusFound, "/")
				}
			}
		}
	}
}

// Google redirects google sign in page
func Google(c *gin.Context) {
	c.Redirect(http.StatusFound, auther.GetGoogleAuthURL())
}

// GoogleCallback registers user info
// @Param state query string true "state"
// @Param code query string true "code"
// @Failure 503 {string} string "no OAuth token"
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if info, err := auther.GetGoogleOAuthInfo(state, code); err != nil {
		c.Set(keyErr, err.Error())
	} else {
		if _, err := services.OAuthSignIn(entities.Google, info); err != nil {
			c.Set(keyErr, err.Error())
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	}
}

// GitHub redirects github sign in page
func GitHub(c *gin.Context) {
	c.Redirect(http.StatusFound, auther.GetGitHubAuthURL())
}

// GitHubCallback registers user info
// @Description try register github info
// @Summary try register github info
// @Param state query string true "state"
// @Param code query string true "code"
// @Failure 503 {string} string "no OAuth token"
func GitHubCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if info, err := auther.GetGitHubOAuthInfo(state, code); err != nil {
		c.HTML(http.StatusServiceUnavailable, "error.tmpl", gin.H{"err": err.Error()})
	} else {
		if _, err := services.OAuthSignIn(entities.GitHub, info); err != nil {
			c.HTML(http.StatusServiceUnavailable, "error.tmpl", gin.H{"err": err.Error()})
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	}
}
