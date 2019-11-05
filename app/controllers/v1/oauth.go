package v1

import (
	"net/http"

	"github.com/gomodule/oauth1/oauth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// Twitter redirects twitter sign in page
// @Description try sign in twitter
// @Summary try sign in twitter
// @Success 302
// @Failure 503 {string} string "failed to get Twitter Request Token"
// @Router /twitter [get]
func Twitter(c *gin.Context) {
	cred, url, err := auther.GetTwitterAuthURL()

	session := sessions.Default(c)
	session.Set("tmpToken", cred.Token)
	session.Set("tmpSecret", cred.Secret)
	session.Save()

	if err != nil {
		c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
	}
	c.Redirect(http.StatusFound, url)
}

// TwitterCallback registers user info
// @Description try register twitter info
// @Summary try register twitter info
// @Param oauth_verifier query string true "access token"
// @Success 302
// @Failure 503 {string} string "session expire or not exists"
// @Router /twitterCallback [get]
func TwitterCallback(c *gin.Context) {
	session := sessions.Default(c)
	secret := c.Query("oauth_verifier")

	if tmpToken := session.Get("tmpToken"); tmpToken == nil {
		c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": "no session key tmpToken"})
	} else {
		if tmpSecret := session.Get("tmpSecret"); tmpSecret == nil {
			c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": "no session key tmpSecret"})
		} else {
			if info, err := auther.GetTwitterOAuthInfo(&oauth.Credentials{
				Token:  tmpToken.(string),
				Secret: tmpSecret.(string),
			}, secret); err != nil {
				c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
			} else {
				if _, err := services.OAuthSignIn(entities.Twitter, info); err != nil {
					c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": "no session key tmpSecret"})
				} else {
					c.Redirect(http.StatusFound, "/")
				}
			}
		}
	}
}

// Google redirects google sign in page
// @Description try sign in google
// @Summary try sign in google
// @Success 302
// @Router /google [get]
func Google(c *gin.Context) {
	c.Redirect(http.StatusFound, auther.GetGoogleAuthURL())
}

// GoogleCallback registers user info
// @Description try register google info
// @Summary try register google info
// @Param state query string true "state"
// @Param code query string true "code"
// @Success 302
// @Failure 503 {string} string "session expire or not exists"
// @Router /googleCallback [get]
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if info, err := auther.GetGoogleOAuthInfo(state, code); err != nil {
		c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
	} else {
		if _, err := services.OAuthSignIn(entities.Google, info); err != nil {
			c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	}
}

// GitHub redirects github sign in page
// @Description try sign in github
// @Summary try sign in github
// @Success 302
// @Router /github [get]
func GitHub(c *gin.Context) {
	c.Redirect(http.StatusFound, auther.GetGitHubAuthURL())
}

// GitHubCallback registers user info
// @Description try register github info
// @Summary try register github info
// @Param state query string true "state"
// @Param code query string true "code"
// @Success 302
// @Failure 503 {string} string "session expire or not exists"
// @Router /githubCallback [get]
func GitHubCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if info, err := auther.GetGitHubOAuthInfo(state, code); err != nil {
		c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
	} else {
		if _, err := services.OAuthSignIn(entities.GitHub, info); err != nil {
			c.HTML(http.StatusServiceUnavailable, "/error.tmpl", gin.H{"err": err.Error()})
		} else {
			c.Redirect(http.StatusFound, "/")
		}
	}
}
