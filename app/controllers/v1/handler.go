package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

var keyErr = "err"
var keyOk = "ok"
var keyOwner = "o"
var keyOwnerErr = "oerr"

// GeneralHandler fetch user from bearer token and set Player to user variable, then respond client
func GeneralHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services.MuModel.Lock()
		defer services.MuModel.Unlock()
		o, err := parseJWT(c.GetHeader("Authorization"))
		c.Set(keyOwner, o)
		c.Set(keyOwnerErr, err)
		c.Next()
		if res, has := c.Get(keyErr); has {
			if verr, ok := res.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, buildErrorMessages(verr))
			} else {
				if e, ok := res.(error); ok {
					c.JSON(http.StatusBadRequest, &errInfo{Err: []string{e.Error()}})
				} else if es, ok := res.([]error); ok {
					var msgs []string
					for _, e := range es {
						msgs = append(msgs, e.Error())
					}
					c.JSON(http.StatusBadRequest, &errInfo{Err: msgs})
				} else {
					c.JSON(http.StatusBadRequest, &errInfo{Err: e})
				}

			}
		} else if res, has := c.Get(keyOk); has {
			c.JSON(http.StatusOK, res)
		} else if err != nil {
			c.JSON(http.StatusUnauthorized, &errInfo{Err: err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, &errInfo{Err: nil})
		}
	}
}

func authorize(c *gin.Context) *entities.Player {
	if o, has := c.Get(keyOwner); has {
		return o.(*entities.Player)
	}
	return nil
}

func parseJWT(header string) (*entities.Player, error) {
	url := conf.Secret.Auth.BaseURL
	token := strings.TrimPrefix(header, "Bearer ")

	obj, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Secret.Auth.Salt), nil
	})

	if err != nil || !obj.Valid {
		return nil, err
	}

	data := obj.Claims.(jwt.MapClaims)
	value := data[fmt.Sprintf("%s/id", url)]
	o, ok := services.Model.Players[uint(value.(float64))]
	if !ok {
		return nil, fmt.Errorf("specified user is already removed")
	}
	return o, nil
}

func buildErrorMessages(errs validator.ValidationErrors) *errInfo {
	msgs := []string{}
	for _, err := range errs {
		if err.Param() == "" {
			msgs = append(msgs, fmt.Sprintf("%s must be %s", err.Field(), err.Tag()))
		} else {
			msgs = append(msgs, fmt.Sprintf("%s must be %s %s", err.Field(), err.Tag(), err.Param()))
		}
	}
	return &errInfo{Err: msgs}
}
