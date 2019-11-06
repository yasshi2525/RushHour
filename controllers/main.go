package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/services"
)

var auther *auth.Auther

// keyErr is set when error occurs
const keyErr = "err"

// keyRedirect is set to redirect page
const keyRedirect = "redirect"

func abortByError(c *gin.Context, err error) {
	c.HTML(http.StatusServiceUnavailable, "error.tmpl", gin.H{"err": err.Error()})
}

// OAuthHandler handles redirect or error page
// It views error page when keyErr is set
// It redirects OAuth sites when keyRedirect is set
func OAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before OAuthHandler")
		defer log.Println("after OAuthHandler")

		c.Next()

		if err, has := c.Get(keyErr); has {
			abortByError(c, err.(error))
		} else {
			c.Redirect(http.StatusFound, c.MustGet(keyRedirect).(string))
		}
	}
}

// CallbackFunc return auth information depending on each service
type CallbackFunc func(string, string) (*auth.OAuthInfo, error)

// keyAuthType is attribute name of entities.AuthType
const keyAuthType = "keyAuthType"

// KeyAuthFunc is attribute name of OAuthFunc
const keyAuthFunc = "keyAuthFunc"

// keyAuthInfo is arrtibute name of fetched data by OAuth
const keyAuthInfo = "keyAuthInfo"

// CallbackHandler parse state and code in order to sign in
func CallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before OAuthHandler")
		defer log.Println("after OAuthHandler")

		c.Next()

		if err, has := c.Get(keyErr); has {
			abortByError(c, err.(error))
		} else {
			state := c.Query("state")
			code := c.Query("code")
			fn := c.MustGet(keyAuthFunc).(CallbackFunc)

			if info, err := fn(state, code); err != nil {
				abortByError(c, err)
			} else {
				c.Set(keyAuthInfo, info)
			}
		}
	}
}

// RegisterHandler sends back jwt token to client
// It should be called after CallbackHandler is called
func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before RegisterHandler")
		defer log.Println("after RegisterHandler")

		c.Next()

		if err, has := c.Get(keyErr); has {
			abortByError(c, err.(error))
		} else {
			ty := c.MustGet(keyAuthType).(entities.AuthType)
			info := c.MustGet(keyAuthInfo).(*auth.OAuthInfo)

			if o, err := services.OAuthSignIn(ty, info); err != nil {
				abortByError(c, err)
			} else if token, err := auther.BuildJWT(o.ExportJWTInfo()); err != nil {
				abortByError(c, err)
			} else {
				c.HTML(http.StatusOK, "oauth.tmpl", gin.H{"jwt": token})
			}
		}
	}
}

// Index returns html containing under maintanance or not
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"inOperation": services.IsInOperation()})
}

// InitController loads config
func InitController(a *auth.Auther) {
	auther = a
}
