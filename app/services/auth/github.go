package auth

import (
	"context"
	"fmt"

	client "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubConf *oauth2.Config

func initGitHub(config cfgOAuth) {
	githubConf = &oauth2.Config{
		ClientID:     config.Client,
		ClientSecret: config.Secret,
		RedirectURL:  fmt.Sprintf("%s/github/callback", baseURL),
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}
}

// GetGitHubAuthURL returns auth url
func GetGitHubAuthURL() string {
	return githubConf.AuthCodeURL(state)
}

// GetGitHubUserInfo returns user info
func GetGitHubUserInfo(resState string, code string) (*UserInfo, error) {
	if resState != state {
		return nil, fmt.Errorf("invalid state")
	}
	ctx := context.Background()
	if token, err := githubConf.Exchange(ctx, code); err != nil {
		return nil, err
	} else if !token.Valid() {
		return nil, fmt.Errorf("invalid token")
	} else {
		c := client.NewClient(githubConf.Client(ctx, token))

		if user, res, err := c.Users.Get(ctx, ""); err != nil {
			return nil, err
		} else {
			defer res.Body.Close()
			if res.StatusCode != 200 {
				buf := make([]byte, 65536)
				res.Body.Read(buf)
				return nil, fmt.Errorf("status %d %+v %s", res.StatusCode, res.Header, buf)
			} else {
				return &UserInfo{
					OAuthToken:  token.AccessToken,
					LoginID:     fmt.Sprintf("%d", *user.ID),
					DisplayName: *user.Name,
					Image:       *user.AvatarURL,
				}, nil
			}
		}
	}
}
