package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/revel/revel"
	"gopkg.in/go-playground/validator.v9"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// loginRequest represents requirement for login
type loginRequest struct {
	// ID is email address of user as login id
	ID string `json:"id" validate:"required,email"`
	// Password is password of binded user
	Password string `json:"password" validate:"required"`
}

// jwtInfo represents authenticate information
type jwtInfo struct {
	// Jwt has json web token representing user information
	Jwt string `json:"jwt"`
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
// @Failure 401 {array} string "reasons of error when login fail"
// @Router /login [post]
func (c API) Login() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()

	params := &loginRequest{}
	if errs := validate.Struct(c.Params.BindJSON(params)); errs != nil {
		c.Response.SetStatus(401)
		return c.RenderJSON(buildErrorMessages(errs.(validator.ValidationErrors)))
	}
	o, err := services.PasswordSignIn(params.ID, params.Password)
	if err != nil {
		c.Response.SetStatus(401)
		return c.RenderJSON([]string{err.Error()})
	}
	return c.RenderJSON(authenticate(o))
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
// @Param hue body integer "player's rail line symbol color (HSV model)"
// @Success 200 {object} jwtInfo "json web token containing user attributes"
// @Failure 422 {array} string "reasons of error when register"
// @Router /register [post]
func (c API) Register() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	params := &registerRequest{}
	if errs := validate.Struct(c.Params.BindJSON(params)); errs != nil {
		c.Response.SetStatus(422)
		return c.RenderJSON(buildErrorMessages(errs.(validator.ValidationErrors)))
	}
	hue, _ := strconv.Atoi(params.Hue)
	o, err := services.PasswordSignUp(params.ID, params.DisplayName, params.Password, hue, entities.Normal)
	if err != nil {
		c.Response.SetStatus(422)
		return c.RenderJSON([]string{err.Error()})
	}
	return c.RenderJSON(authenticate(o))
}

// GetSettings returns the list of customizable attributes
// @Description list up user attributes including private one
// @Tags services.AccountSettings
// @Summary get user attributes
// @Accept json
// @Produce json
// @Param Authorization header string true "with the bearer started"
// @Success 200 {object} services.AccountSettings "user attributes"
// @Failure 401 {array} string "invalid jwt"
// @Router /settings [get]
func (c API) GetSettings() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()

	o, err := parse(c.Request.GetHttpHeader("Authorization"))
	if err != nil {
		c.Response.SetStatus(401)
		return c.RenderJSON([]string{err.Error()})
	}
	return c.RenderJSON(services.GetAccountSettings(o.O))
}

// settingsCName is format of custom user name
type settingsCName struct {
	Value string `json:"value" validation:"required"`
}

type settingsUseCname struct {
	Value string `json:"value" validation:"required,eq=true|eq=false"`
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
// @Failure 401 {array} string "invalid jwt"
// @Failure 422 {array} string "invalid parameter"
// @Router /settings/{resname} [post]
func (c API) ChangeSettings() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	o, err := parse(c.Request.GetHttpHeader("Authorization"))
	if err != nil {
		c.Response.SetStatus(401)
		return c.RenderJSON([]string{err.Error()})
	}

	res := c.Params.Get("resname")
	switch res {
	case "custom_name":
		val := &settingsCName{}
		if errs := validate.Struct(c.Params.BindJSON(val)); err != nil {
			c.Response.SetStatus(422)
			return c.RenderJSON(buildErrorMessages(errs.(validator.ValidationErrors)))
		}
		o.O.CustomDisplayName = auth.Encrypt(val.Value)
		return c.RenderJSON(entry{Key: res, Value: val.Value})
	case "use_cname":
		val := &settingsUseCname{}
		if errs := validate.Struct(c.Params.BindJSON(val)); err != nil {
			c.Response.SetStatus(422)
			return c.RenderJSON(buildErrorMessages(errs.(validator.ValidationErrors)))
		}
		o.O.UseCustomDisplayName, _ = strconv.ParseBool(val.Value)
		return c.RenderJSON(entry{Key: res, Value: val.Value})
	default:
		return c.RenderJSON([]error{fmt.Errorf("invalid attribute %s", res)})
	}
}

// SignOut deletes cached OAuth token.
// @Description deletes OAuth token
// @Summary execute user sign out
// @Accept json
// @Produce json
// @Success 200 {object} null "sign out successfully"
// @Failure 401 {array} string "invalid jwt"
// @Router /settings/{resname} [post]
func (c API) SignOut() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	o, err := parse(c.Request.GetHttpHeader("Authorization"))
	if err != nil {
		c.Response.SetStatus(401)
		return c.RenderJSON([]string{err.Error()})
	}
	services.SignOut(o)
	return c.RenderJSON(nil)
}

func authenticate(o *entities.Player) *jwtInfo {
	url := services.Secret.Auth.BaseURL
	now := time.Now()
	exp := now.Add(time.Hour)
	uu := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":                        url,
		"sub":                        "AccessToken",
		"aud":                        url,
		"exp":                        exp.Unix(),
		"nbf":                        now.Unix(),
		"iat":                        now.Unix(),
		"jti":                        uu.String(),
		fmt.Sprintf("%s/id", url):    o.ID,
		fmt.Sprintf("%s/name", url):  auth.Decrypt(o.GetDisplayName()),
		fmt.Sprintf("%s/image", url): auth.Decrypt(o.GetImage()),
		fmt.Sprintf("%s/admin", url): o.Level == entities.Admin,
		fmt.Sprintf("%s/hue", url):   o.Hue,
	})

	jwt, err := token.SignedString([]byte(services.Secret.Auth.Salt))
	if err != nil {
		panic(err)
	}
	return &jwtInfo{jwt}
}

func parse(header string) (*entities.Player, error) {
	url := services.Secret.Auth.BaseURL
	token := strings.TrimPrefix(header, "Bearer ")
	if obj, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(services.Secret.Auth.Salt), nil
	}); err != nil || !obj.Valid {
		return nil, err
	} else {
		data := obj.Claims.(jwt.MapClaims)
		value := data[fmt.Sprintf("%s/id", url)]
		o, ok := services.Model.Players[uint(value.(float64))]
		if !ok {
			return nil, fmt.Errorf("specified user is already removed")
		}
		return o, nil
	}
}
