package auth

import (
	"testing"

	"github.com/yasshi2525/RushHour/app/config"
)

const plainEmpty = ""
const encEmpty = "H24L0VncY1UJva7/H5t0BdE40Qm/8yrq1wqhQvd8oQk="
const digestEmpty = "uTbO6Gyfh6pdPG8uhMtaQjml/lBICm7Ga3CrWx9KxnMMbFFUIbMn7B1pQC5T37Sa1zgesGezOP17DLIiRyJdRw=="

const plainValue = "non-zero"
const encValue = "4Y3XEKaZrYpB5djETs9cM9DnOeLmuDxr2IKupYKXmtM="
const digestValue = "AOChf1F7i9DWMvrd8EXUYbv5V9FlXU+HbBv8DIM4X5hLz+8U48Fy2z8Y6aGJndEj3W7YeQwLd+xwE36l1DwQWA=="

func TestGetAuther(t *testing.T) {
	if _, err := GetAuther(config.CnfAuth{}); err != nil {
		t.Errorf("GetAuther({}) got %v, want nil", err)
	}
}

func TestOAuth(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in   *OAuthInfo
			want bool
		}{
			{
				in: &OAuthInfo{
					Handler:    a,
					IsEnc:      false,
					OAuthToken: plainEmpty,
				},
				want: false,
			}, {
				in: &OAuthInfo{
					Handler:    a,
					IsEnc:      false,
					OAuthToken: plainValue,
				},
				want: true,
			}, {
				in: &OAuthInfo{
					Handler:    a,
					IsEnc:      true,
					OAuthToken: encEmpty,
				},
				want: false,
			}, {
				in: &OAuthInfo{
					Handler:    a,
					IsEnc:      true,
					OAuthToken: encValue,
				},
				want: true,
			},
		}

		for _, c := range cases {
			if got := c.in.IsValid(); got != c.want {
				t.Errorf("%v.IsValid() got %t, want %t", c.in, got, c.want)
			}
		}
	})

	t.Run("Enc.err", func(t *testing.T) {
		i := &OAuthInfo{IsEnc: true}
		if _, err := i.Enc(); err == nil {
			t.Errorf("%v.Enc() got nil, want not nil", i)
		}
	})

	t.Run("Enc", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in      *OAuthInfo
			wantNot string
		}{
			{
				in: &OAuthInfo{
					Handler:     a,
					IsEnc:       false,
					DisplayName: plainValue,
					Image:       plainValue,
					LoginID:     plainValue,
					OAuthToken:  plainValue,
					OAuthSecret: plainValue,
				},
				wantNot: plainValue,
			},
		}

		for _, c := range cases {
			got, err := c.in.Enc()
			if err != nil {
				t.Errorf("%v.Enc().err got %v, want nil", c.in, err)
			}
			if got.DisplayName == c.wantNot {
				t.Errorf("%v.Enc().DisplayName got %s, want not %s", c.in, got.DisplayName, c.wantNot)
			}
			if got.Image == c.wantNot {
				t.Errorf("%v.Enc().Image got %s, want not %s", c.in, got.Image, c.wantNot)
			}
			if got.LoginID == c.wantNot {
				t.Errorf("%v.Enc().LoginID got %s, want not %s", c.in, got.LoginID, c.wantNot)
			}
			if got.OAuthToken == c.wantNot {
				t.Errorf("%v.Enc().OAuthToken got %s, want not %s", c.in, got.OAuthToken, c.wantNot)
			}
			if got.OAuthSecret == c.wantNot {
				t.Errorf("%v.Enc().OAuthSecret got %s, want not %s", c.in, got.OAuthSecret, c.wantNot)
			}
		}
	})

	t.Run("Dec.err", func(t *testing.T) {
		i := &OAuthInfo{}
		if _, err := i.Dec(); err == nil {
			t.Errorf("%v.Dec.err got nil, want not nil", i)
		}
	})

	t.Run("Dec", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in   *OAuthInfo
			want string
		}{
			{
				in: &OAuthInfo{
					Handler:     a,
					IsEnc:       true,
					DisplayName: encValue,
					Image:       encValue,
					LoginID:     encValue,
					OAuthToken:  encValue,
					OAuthSecret: encValue,
				},
				want: plainValue,
			},
		}

		for _, c := range cases {
			got, err := c.in.Dec()
			if err != nil {
				t.Errorf("%v.Dec().err got %v, want nil", c.in, err)
			}
			if got.DisplayName != c.want {
				t.Errorf("%v.Dec().DisplayName got %s, want %s", c.in, got.DisplayName, c.want)
			}
			if got.Image != c.want {
				t.Errorf("%v.Dec().Image got %s, want %s", c.in, got.Image, c.want)
			}
			if got.LoginID != c.want {
				t.Errorf("%v.Dec().LoginID got %s, want %s", c.in, got.LoginID, c.want)
			}
			if got.OAuthToken != c.want {
				t.Errorf("%v.Dec().OAuthToken got %s, want %s", c.in, got.OAuthToken, c.want)
			}
			if got.OAuthSecret != c.want {
				t.Errorf("%v.Dec().OAuthSecret got %s, want %s", c.in, got.OAuthSecret, c.want)
			}
		}
	})
}

func TestAuther(t *testing.T) {
	t.Run("Digest", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in   string
			want string
		}{
			{
				in:   plainEmpty,
				want: digestEmpty,
			},
			{
				in:   plainValue,
				want: digestValue,
			},
		}

		for _, c := range cases {
			if got := a.Digest(c.in); got != c.want {
				t.Errorf("Digest(%s) got %s, want %s", c.in, got, c.want)
			}
		}
	})
	t.Run("Encrypt", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in      string
			wantNot string
		}{
			{
				in:      plainEmpty,
				wantNot: plainEmpty,
			},
			{
				in:      plainValue,
				wantNot: plainEmpty,
			},
		}

		for _, c := range cases {
			if got := a.Encrypt(c.in); got == c.wantNot {
				t.Errorf("Encrypt(%s) got %s, want not %s", c.in, got, c.wantNot)
			}
		}
	})
	t.Run("Decrypt", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		cases := []struct {
			in   string
			want string
		}{
			{
				in:   plainEmpty,
				want: plainEmpty,
			},
			{
				in:   encEmpty,
				want: plainEmpty,
			},
			{
				in:   encValue,
				want: plainValue,
			},
		}

		for _, c := range cases {
			if got := a.Decrypt(c.in); got != c.want {
				t.Errorf("Decrypt(%s) got %s, want %s", c.in, got, c.want)
			}
		}
	})
	t.Run("BuildJWT", func(t *testing.T) {
		a, _ := GetAuther(config.CnfAuth{})
		if _, err := a.BuildJWT(&JWTInfo{}); err != nil {
			t.Errorf("buildJWT().err got %v, want nil", err)
		}
	})
}
