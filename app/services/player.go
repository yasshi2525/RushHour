package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string, level entities.PlayerType) (*entities.Player, error) {
	if o, err := Model.PasswordSignUp(loginid, password); err != nil {
		return nil, err
	} else {
		o.Level = level
		o.DisplayName = displayname
		url := fmt.Sprintf("%s/public/img/player.png", Secret.Auth.BaseURL)
		o.Image = auth.Encrypt(url)
		AddOpLog("CreatePlayer", o)
		return o, nil
	}
}

// OAuthSignIn find or create Player by OAuth
func OAuthSignIn(authType entities.AuthType, info *auth.UserInfo) (*entities.Player, error) {
	return Model.OAuthSignIn(authType, info)
}

// SignOut delete Player's token value
func SignOut(token string) {
	if o, found := Model.Tokens[token]; found {
		o.SignOut()
	}
}

// PasswordSignIn finds Player by loginid and password
func PasswordSignIn(loginid string, password string) (*entities.Player, error) {
	return Model.PasswordSignIn(loginid, password)
}

// PasswordSignUp creates Player with loginid and password
func PasswordSignUp(loginid string, password string) (*entities.Player, error) {
	return Model.PasswordSignUp(loginid, password)
}

// FindOwner returns Player by token
func FindOwner(token string) *entities.Player {
	return Model.Tokens[token]
}
