package v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

func authenticate(o *entities.Player) string {
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
	return jwt
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
		return services.Model.Players[uint(value.(float64))], nil
	}
}
