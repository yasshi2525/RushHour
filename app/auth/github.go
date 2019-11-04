package auth

import (
	"context"
	"fmt"

	client "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/yasshi2525/RushHour/app/config"
)

func (a *Auther) initGitHub(conf config.CnfOAuth) {
	a.githubConf = &oauth2.Config{
		ClientID:     conf.Client,
		ClientSecret: conf.Secret,
		RedirectURL:  fmt.Sprintf("%s/github/callback", a.baseURL),
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}
}

// GetGitHubAuthURL returns auth url
func (a *Auther) GetGitHubAuthURL() string {
	return a.githubConf.AuthCodeURL(a.state)
}

// GetGitHubOAuthInfo returns user info
func (a *Auther) GetGitHubOAuthInfo(resState string, code string) (*OAuthInfo, error) {
	if resState != a.state {
		return nil, fmt.Errorf("invalid state")
	}
	ctx := context.Background()
	if token, err := a.githubConf.Exchange(ctx, code); err != nil {
		return nil, err
	} else if !token.Valid() {
		return nil, fmt.Errorf("invalid token")
	} else {
		c := client.NewClient(a.githubConf.Client(ctx, token))
		user, res, err := c.Users.Get(ctx, "")
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			buf := make([]byte, 65536)
			res.Body.Read(buf)
			return nil, fmt.Errorf("status %d %+v %s", res.StatusCode, res.Header, buf)
		}
		return &OAuthInfo{
			handler:     a,
			OAuthToken:  token.AccessToken,
			LoginID:     fmt.Sprintf("%d", *user.ID),
			DisplayName: *user.Name,
			Image:       *user.AvatarURL,
		}, nil
	}
}
