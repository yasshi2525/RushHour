package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gomodule/oauth1/oauth"

	"github.com/yasshi2525/RushHour/app/config"
)

var twitterClient *oauth.Client

func (a *Auther) initTwitter(conf config.CnfTwitter) {
	twitterClient = &oauth.Client{
		TemporaryCredentialRequestURI: conf.Request,
		ResourceOwnerAuthorizationURI: conf.Authenticate,
		TokenRequestURI:               conf.AccessToken,
		Credentials: oauth.Credentials{
			Token:  conf.Token,
			Secret: conf.Secret,
		},
	}
}

// GetTwitterAuthURL returns auth url
func (a *Auther) GetTwitterAuthURL() (*oauth.Credentials, string) {
	tmpCred, err := twitterClient.RequestTemporaryCredentials(
		&http.Client{}, fmt.Sprintf("%s/twitter/callback", a.baseURL), nil)

	if err != nil {
		panic(err)
	}

	return tmpCred, twitterClient.AuthorizationURL(tmpCred, nil)
}

type twitterOAuthInfo struct {
	LoginID     string `json:"id_str"`
	DisplayName string `json:"name"`
	Image       string `json:"profile_image_url_https"`
}

// GetTwitterOAuthInfo returns user info
func (a *Auther) GetTwitterOAuthInfo(tmpCred *oauth.Credentials, tmpSecret string) (*OAuthInfo, error) {
	cred, values, err := twitterClient.RequestToken(&http.Client{}, tmpCred, tmpSecret)
	if err != nil {
		return nil, err
	}
	res, err := twitterClient.Get(&http.Client{}, cred, "https://api.twitter.com/1.1/users/show.json", url.Values{"user_id": {values.Get("user_id")}})
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		buf := make([]byte, 65536)
		res.Body.Read(buf)
		return nil, fmt.Errorf("status %d %+v %s", res.StatusCode, res.Header, buf)
	}

	info := &twitterOAuthInfo{}
	if err := json.NewDecoder(res.Body).Decode(info); err != nil {
		return nil, err
	}
	return &OAuthInfo{
		OAuthToken:  values.Get("oauth_token"),
		OAuthSecret: values.Get("oauth_token_secret"),
		LoginID:     info.LoginID,
		DisplayName: info.DisplayName,
		Image:       info.Image,
	}, nil
}
