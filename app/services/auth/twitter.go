package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gomodule/oauth1/oauth"
)

var twitterClient *oauth.Client

func initTwitter(config cfgTwitter) {
	twitterClient = &oauth.Client{
		TemporaryCredentialRequestURI: config.Request,
		ResourceOwnerAuthorizationURI: config.Authenticate,
		TokenRequestURI:               config.AccessToken,
		Credentials: oauth.Credentials{
			Token:  config.Token,
			Secret: config.Secret,
		},
	}
}

// GetTwitterAuthURL returns auth url
func GetTwitterAuthURL() (*oauth.Credentials, string) {
	tmpCred, err := twitterClient.RequestTemporaryCredentials(
		&http.Client{}, fmt.Sprintf("%s/twitter/callback", baseURL), nil)

	if err != nil {
		panic(err)
	}

	return tmpCred, twitterClient.AuthorizationURL(tmpCred, nil)
}

type twitterUserInfo struct {
	LoginID     string `json:"id_str"`
	DisplayName string `json:"name"`
	Image       string `json:"profile_image_url_https"`
}

// GetTwitterUserInfo returns user info
func GetTwitterUserInfo(tmpCred *oauth.Credentials, tmpSecret string) (*UserInfo, error) {
	if cred, values, err := twitterClient.RequestToken(&http.Client{}, tmpCred, tmpSecret); err != nil {
		return nil, err
	} else {
		if res, err := twitterClient.Get(&http.Client{}, cred,
			"https://api.twitter.com/1.1/users/show.json", url.Values{"user_id": {values.Get("user_id")}}); err != nil {
			return nil, err
		} else {
			defer res.Body.Close()
			if res.StatusCode != 200 {
				buf := make([]byte, 65536)
				res.Body.Read(buf)
				return nil, fmt.Errorf("status %d %+v %s", res.StatusCode, res.Header, buf)
			} else {
				info := &twitterUserInfo{}
				if err := json.NewDecoder(res.Body).Decode(info); err != nil {
					return nil, err
				} else {
					return &UserInfo{
						OAuthToken:  values.Get("oauth_token"),
						OAuthSecret: values.Get("oauth_token_secret"),
						LoginID:     info.LoginID,
						DisplayName: info.DisplayName,
						Image:       info.Image,
					}, nil
				}
			}
		}
	}
}
