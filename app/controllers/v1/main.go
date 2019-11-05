package v1

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/yasshi2525/RushHour/app/auth"
	"github.com/yasshi2525/RushHour/app/config"
	"github.com/yasshi2525/RushHour/app/entities"
)

// entry represents generic key-value pair
type entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type jwtInfo struct {
	Jwt string `json:"jwt"`
}

// user represents public attributes that everyone can view
type user struct {
	// ID is number
	ID uint `json:"id"`
	// Name is display name
	Name string `json:"name"`
	// Image is url of profile icon
	Image string `json:"image"`
	// Hue is rail line color (HSV model)
	Hue float64 `json:"hue"`
}

type errInfo struct {
	Err interface{} `json:"err"`
}

var conf config.Config
var auther *auth.Auther

// buildJwt returns JSON Web Token of player
func buildJwt(o *entities.Player) (*jwtInfo, error) {
	url := conf.Secret.Auth.BaseURL
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
		fmt.Sprintf("%s/name", url):  auther.Decrypt(o.GetDisplayName()),
		fmt.Sprintf("%s/image", url): auther.Decrypt(o.GetImage()),
		fmt.Sprintf("%s/admin", url): o.Level == entities.Admin,
		fmt.Sprintf("%s/hue", url):   o.Hue,
	})

	jwt, err := token.SignedString([]byte(conf.Secret.Auth.Salt))
	if err != nil {
		return nil, err
	}
	return &jwtInfo{jwt}, nil
}

// InitController loads config
func InitController(c config.Config, a *auth.Auther) {
	conf = c
	auther = a
}
