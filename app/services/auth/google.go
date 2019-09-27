package auth

import (
	"context"
	"fmt"

	"google.golang.org/api/option"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	v2 "google.golang.org/api/oauth2/v2"
)

var googleConf *oauth2.Config

func initGoogle(config cfgOAuth) {
	googleConf = &oauth2.Config{
		ClientID:     config.Client,
		ClientSecret: config.Secret,
		RedirectURL:  fmt.Sprintf("%s/google/callback", baseURL),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email", "profile"},
	}
}

// GetGoogleAuthURL returns auth url
func GetGoogleAuthURL() string {
	return googleConf.AuthCodeURL(state)
}

// GetGoogleUserInfo returns user info
func GetGoogleUserInfo(resState string, code string) (*UserInfo, error) {
	if resState != state {
		return nil, fmt.Errorf("invalid state")
	}
	ctx := context.Background()
	if token, err := googleConf.Exchange(ctx, code); err != nil {
		return nil, err
	} else if !token.Valid() {
		return nil, fmt.Errorf("invalid token")
	} else {
		if service, err := v2.NewService(ctx, option.WithTokenSource(googleConf.TokenSource(ctx, token))); err != nil {
			return nil, err
		} else {
			if info, err := service.Tokeninfo().Do(); err != nil {
				return nil, err
			} else {
				if person, err := service.Userinfo.V2.Me.Get().Do(); err != nil {
					return nil, err
				} else {
					return &UserInfo{
						OAuthToken:  token.AccessToken,
						LoginID:     info.UserId,
						DisplayName: person.Name,
						Image:       person.Picture,
					}, nil
				}
			}
		}
	}
}
