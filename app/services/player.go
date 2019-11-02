package services

import (
	"encoding/json"
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/auth"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string, hue int, lv entities.PlayerType) (*entities.Player, error) {
	if err := CheckMaintenance(); err != nil {
		return nil, err
	}
	if o, err := Model.PasswordSignUp(loginid, password, lv); err != nil {
		return nil, err
	} else {
		o.CustomDisplayName = auth.Encrypt(displayname)
		o.UseCustomDisplayName = true
		o.Hue = hue
		url := fmt.Sprintf("%s/public/img/player.png", Secret.Auth.BaseURL)
		o.CustomImage = auth.Encrypt(url)
		o.UseCustomImage = true
		AddOpLog("CreatePlayer", o)
		return o, nil
	}
}

// OAuthSignIn find or create Player by OAuth
func OAuthSignIn(authType entities.AuthType, info *auth.OAuthInfo) (*entities.Player, error) {
	if o, err := Model.OAuthSignIn(authType, info); err != nil {
		return nil, err
	} else if err := CheckMaintenance(o); err != nil {
		return nil, err
	} else {
		return o, nil
	}
}

// SignOut delete Player's token value
func SignOut(token string) {
	if o, found := Model.Tokens[token]; found {
		o.SignOut()
	}
}

// PasswordSignIn finds Player by loginid and password
func PasswordSignIn(loginid string, password string) (*entities.Player, error) {
	if o, err := Model.PasswordSignIn(loginid, password); err != nil {
		return nil, err
	} else if err := CheckMaintenance(o); err != nil {
		return nil, err
	} else {
		return o, nil
	}
}

// PasswordSignUp creates Player with loginid and password
func PasswordSignUp(loginid string, name string, password string, hue int, lv entities.PlayerType) (*entities.Player, error) {
	if err := CheckMaintenance(); lv != entities.Admin && err != nil {
		return nil, err
	}
	if o, err := Model.PasswordSignUp(loginid, password, lv); err != nil {
		return nil, err
	} else {
		o.CustomDisplayName = auth.Encrypt(name)
		o.UseCustomDisplayName = true
		o.Hue = hue
		url := fmt.Sprintf("%s/public/img/player.png", Secret.Auth.BaseURL)
		o.CustomImage = auth.Encrypt(url)
		o.UseCustomImage = true
		return o, nil
	}
}

// AccountSettings returns user customizable attributes.
type AccountSettings struct {
	Player      *entities.Player  `json:"-"`
	CustomName  string            `json:"custom_name"`
	CustomImage string            `json:"custom_image"`
	AuthType    entities.AuthType `json:"auth_type"`
}

// MarshalJSON returns plain text data.
func (s *AccountSettings) MarshalJSON() ([]byte, error) {
	type Alias AccountSettings
	if s.Player.Auth == entities.Local {
		return json.Marshal(&struct {
			LoginID string `json:"email"`
			*Alias
		}{
			LoginID: auth.Decrypt(s.Player.LoginID),
			Alias:   (*Alias)(s),
		})
	}
	return json.Marshal(&struct {
		OAuthName      string `json:"oauth_name"`
		UseCustomName  bool   `json:"use_cname"`
		OAuthImage     string `json:"oauth_image"`
		UseCustomImage bool   `json:"use_cimage"`
		*Alias
	}{
		OAuthName:      auth.Decrypt(s.Player.OAuthDisplayName),
		UseCustomName:  s.Player.UseCustomDisplayName,
		OAuthImage:     auth.Decrypt(s.Player.OAuthImage),
		UseCustomImage: s.Player.UseCustomImage,
		Alias:          (*Alias)(s),
	})
}

// GetAccountSettings returns the list of customizable attributes.
func GetAccountSettings(o *entities.Player) *AccountSettings {
	return &AccountSettings{
		Player:      o,
		CustomName:  auth.Decrypt(o.CustomDisplayName),
		CustomImage: auth.Decrypt(o.CustomImage),
		AuthType:    o.Auth,
	}
}
