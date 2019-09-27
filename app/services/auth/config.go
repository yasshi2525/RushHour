package auth

import "crypto/cipher"

var baseURL string
var state string
var salt string
var block cipher.Block

// UserInfo represents infomation from OAuth App
type UserInfo struct {
	DisplayName string
	Image       string
	LoginID     string
	OAuthToken  string
	OAuthSecret string
	IsEnc       bool
}

type cfgTwitter struct {
	Token        string
	Secret       string
	Request      string `validate:"url"`
	Authenticate string `validate:"url"`
	AccessToken  string `validate:"url"`
}

type cfgOAuth struct {
	Client string
	Secret string
}

// Config represents secret information structure.
type Config struct {
	BaseURL string `validate:"url"`
	Salt    string
	Key     string `validate:"len=16"`
	State   string

	Twitter cfgTwitter
	Google  cfgOAuth
}
