package services

import (
	"fmt"
	"strings"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string, level entities.PlayerType) (*entities.Player, error) {
	// duplication check
	for _, oth := range Model.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			revel.AppLog.Warnf("Login ID \"%s\" is already used.", loginid)
			return nil, fmt.Errorf("login ID \"%s\" is already used", loginid)
		}
	}

	player := Model.NewPlayer()
	player.LoginID = loginid
	player.DisplayName = displayname
	player.Password = password
	player.Level = level

	AddOpLog("CreatePlayer", player)
	return player, nil
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
	return Model.PasswordSignIn(auth.Digest(loginid), auth.Encrypt(password))
}

// PasswordSignUp creates Player with loginid and password
func PasswordSignUp(loginid string, password string) (*entities.Player, error) {
	return Model.PasswordSignUp(auth.Digest(loginid), auth.Encrypt(password))
}

// FetchOwner fetch loginid Player.
func FetchOwner(loginid string) (*entities.Player, error) {
	for _, oth := range Model.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			return oth, nil
		}
	}
	return nil, fmt.Errorf("login ID \"%s\" was not found", loginid)
}

// FindOwner returns Player by token
func FindOwner(token string) *entities.Player {
	return Model.Tokens[token]
}
