package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	v2 "google.golang.org/api/oauth2/v2"

	"github.com/yasshi2525/RushHour/app/config"
)

func (a *Auther) initGoogle(conf config.CnfOAuth) {
	a.googleConf = &oauth2.Config{
		ClientID:     conf.Client,
		ClientSecret: conf.Secret,
		RedirectURL:  fmt.Sprintf("%s/google/callback", a.baseURL),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email", "profile"},
	}
}

// GetGoogleAuthURL returns auth url
func (a *Auther) GetGoogleAuthURL() string {
	return a.googleConf.AuthCodeURL(a.state)
}

// GetGoogleOAuthInfo returns user info
func (a *Auther) GetGoogleOAuthInfo(resState string, code string) (*OAuthInfo, error) {
	if resState != a.state {
		return nil, fmt.Errorf("invalid state")
	}
	ctx := context.Background()
	if token, err := a.googleConf.Exchange(ctx, code); err != nil {
		return nil, err
	} else if !token.Valid() {
		return nil, fmt.Errorf("invalid token")
	} else {
		service, err := v2.NewService(ctx, option.WithTokenSource(a.googleConf.TokenSource(ctx, token)))
		if err != nil {
			return nil, err
		}
		info, err := service.Tokeninfo().Do()
		if err != nil {
			return nil, err
		}
		person, err := service.Userinfo.V2.Me.Get().Do()
		if err != nil {
			return nil, err
		}
		return &OAuthInfo{
			handler:     a,
			OAuthToken:  token.AccessToken,
			LoginID:     info.UserId,
			DisplayName: person.Name,
			Image:       person.Picture,
		}, nil
	}
}
