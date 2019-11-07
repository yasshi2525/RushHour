package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/services"
)

// loginRequest represents requirement for login
type loginRequest struct {
	// ID is email address of user as login id
	ID string `form:"id" json:"id" validate:"required,email"`
	// Password is password of binded user
	Password string `form:"password" json:"password" validate:"required"`
}

// Login returns result of password login
// @Description try login using loginid/password paramter
// @Tags jwtInfo
// @Summary try login to RushHour server
// @Accept json
// @Produce json
// @Param id body string true "email address"
// @Param password body string true "password"
// @Success 200 {object} jwtInfo "json web token containing user attributes"
// @Failure 400 {object} errInfo "reasons of error when login fail"
// @Failure 503 {object} errInfo "under maintenance (apply only normal, except admin)"
// @Router /login [post]
func Login(c *gin.Context) {
	params := loginRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else {
		if o, err := services.PasswordSignIn(params.ID, params.Password); err != nil {
			c.Set(keyErr, err)
		} else if !services.IsInOperation() && o.Level == entities.Normal {
			abortByMaintenance(c)
		} else if jwt, err := auther.BuildJWT(o.ExportJWTInfo()); err != nil {
			c.Set(keyErr, err)
		} else {
			c.Set(keyOk, &jwtInfo{jwt})
		}
	}
}

// registerRequest represents requirement for sign up to RushHour server
type registerRequest struct {
	loginRequest
	// DisplayName is shown by everyone (default: NoName)
	DisplayName string `json:"name" validate:"omitempty"`
	// Hue is rail line symbol color (HSV model)
	Hue string `json:"hue" validate:"required,numeric"`
}

// validRegisterRequest validates that Register satisfies sign in conditions
func validRegisterRequest(sl validator.StructLevel) {
	v := sl.Current().Interface().(registerRequest)
	hue, _ := strconv.Atoi(v.Hue)

	if hue < 0 {
		sl.ReportError(v.Hue, "hue", "Hue", "gte", "0")
	}
	if hue >= 360 {
		sl.ReportError(v.Hue, "hue", "Hue", "lt", "360")
	}
}

// Register returns result of password sign up
// @Description try register with loginid/password
// @Tags jwtInfo
// @Summary try register to RushHour server
// @Accept json
// @Produce json
// @Param id body string true "email address"
// @Param password body string true "password"
// @Param name body string false "display name"
// @Param hue body integer true "player's rail line symbol color (HSV model)"
// @Success 200 {object} jwtInfo "json web token containing user attributes"
// @Failure 400 {object} errInfo "reasons of error when register"
// @Failure 503 {object} errInfo "under maintenance"
// @Router /register [post]
func Register(c *gin.Context) {
	params := registerRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else {
		hue, _ := strconv.Atoi(params.Hue)
		if o, err := services.PasswordSignUp(params.ID, params.DisplayName, params.Password, hue, entities.Normal); err != nil {
			c.Set(keyErr, err)
		} else if jwt, err := auther.BuildJWT(o.ExportJWTInfo()); err != nil {
			c.Set(keyErr, err)
		} else {
			c.Set(keyOk, &jwtInfo{jwt})
		}
	}
}

// Settings returns the list of customizable attributes
// @Description list up user attributes including private one
// @Tags services.AccountSettings
// @Summary get user attributes
// @Accept json
// @Produce json
// @Param Authorization header string true "with the bearer started"
// @Success 200 {object} services.AccountSettings "user attributes"
// @Failure 400 {object} errInfo "under maintenance"
// @Failure 401 {object} errInfo "invalid jwt"
// @Failure 503 {object} errInfo "under maintenance"
// @Router /settings [get]
func Settings(c *gin.Context) {
	o := c.MustGet(keyOwner).(*entities.Player)
	c.Set(keyOk, services.GetAccountSettings(o))
}

// settingsCName is format of custom user name
type settingsCName struct {
	Value string `form:"value" json:"value" validation:"required"`
}

type settingsUseCname struct {
	Value string `form:"value" json:"value" validation:"required,eq=true|eq=false"`
}

// ChangeSettings returns the result of change profile
// @Description change user attributes including private one
// @Tags entry
// @Summary change user attributes
// @Accept json
// @Produce json
// @Param Authorization header string true "with the bearer started"
// @Param resname path string true "changing resource name" Enums(custom_name, use_cname)
// @Param value body string true "changing resource value"
// @Success 200 {object} entry "changed user attributes"
// @Failure 400 {object} errInfo "invalid parameter"
// @Failure 401 {object} errInfo "invalid jwt"
// @Failure 503 {object} errInfo "under maintenance"
// @Router /settings/{resname} [post]
func ChangeSettings(c *gin.Context) {
	o := c.MustGet(keyOwner).(*entities.Player)
	res := c.Param("resname")
	switch res {
	case "custom_name":
		val := settingsCName{}
		if err := c.Bind(&val); err != nil {
			c.Set(keyErr, err)
		} else {
			o.CustomDisplayName = auther.Encrypt(val.Value)
			c.Set(keyOk, entry{Key: res, Value: val.Value})
		}
	case "use_cname":
		val := settingsUseCname{}
		if err := c.Bind(&val); err != nil {
			c.Set(keyErr, err)
		} else {
			o.UseCustomDisplayName, _ = strconv.ParseBool(val.Value)
			c.Set(keyOk, entry{Key: res, Value: val.Value})
		}
	default:
		c.Set(keyErr, fmt.Errorf("invalid attribute %s", res))
	}
}

// SignOut deletes cached OAuth token.
// @Description deletes OAuth token
// @Summary execute user sign out
// @Accept json
// @Produce json
// @Success 200 {object} string "sign out successfully"
// @Failure 401 {object} errInfo "invalid jwt"
// @Failure 503 {object} errInfo "under maintenance"
// @Router /signout [get]
func SignOut(c *gin.Context) {
	o := c.MustGet(keyOwner).(*entities.Player)
	services.SignOut(o)
	c.Set(keyOk, nil)
}
