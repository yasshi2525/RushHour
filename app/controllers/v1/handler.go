package v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// keyErr is set when error happens
const keyErr = "err"

// keyOk is set when no error happens
const keyOk = "ok"

// keyOwner is set when user is specified by jwt
const keyOwner = "o"

// keyOAuth is set when user information is received by OAuth
const keyOAuth = "oauth"

func abortByMaintenance(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, &errInfo{Err: []string{"under maintenance"}})
	c.Abort()
}

// MaintenanceHandler blocks user access under maintenance
func MaintenanceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before MaintenanceHandler")
		defer log.Println("after MaintenanceHandler")
		if !services.IsInOperation() {
			abortByMaintenance(c)
		} else {
			c.Next()
		}
	}
}

// JWTHandler handles user action with jwt key
// It should be called after MaintenanceHandler is called
func JWTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before JWTHandler")
		defer log.Println("after JWTHandler")
		o, err := parseJWT(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, &errInfo{Err: []string{err.Error()}})
			c.Abort()
		}
		c.Set(keyOwner, o)
		c.Next()
	}
}

// AdminHandler handles admin action
// It should be called after JWTHandler is called
func AdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("before AdminHandler")
		defer log.Println("after AdminHandler")
		if c.MustGet(keyOwner).(*entities.Player).Level != entities.Admin {
			c.JSON(http.StatusUnauthorized, &errInfo{Err: []string{"permission denied"}})
			c.Abort()
		}
		c.Next()
	}
}

// ModelHandler handles model and error controling
// It causes panic when neither result keyOk or keyErr is set
// It should be called after JWTHandeler and AdminHandler are called
func ModelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services.MuModel.Lock()
		defer services.MuModel.Unlock()
		log.Println("before ModelHandler")
		defer log.Println("after ModelHandler")
		c.Next()
		// error reported
		if res, has := c.Get(keyErr); has {
			// error caused by validation
			if verr, ok := res.(validator.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, buildErrorMessages(verr))
			} else {
				// error caused by services
				if e, ok := res.(error); ok {
					// single reason
					c.JSON(http.StatusBadRequest, &errInfo{Err: []string{e.Error()}})
				} else if es, ok := res.([]error); ok {
					// multiple reason
					var msgs []string
					for _, e := range es {
						msgs = append(msgs, e.Error())
					}
					c.JSON(http.StatusBadRequest, &errInfo{Err: msgs})
				} else {
					// unhandle error
					c.JSON(http.StatusBadRequest, &errInfo{Err: []string{fmt.Sprintf("%s", e)}})
				}

			}
		} else {
			c.JSON(http.StatusOK, c.MustGet(keyOk))
		}
	}
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
