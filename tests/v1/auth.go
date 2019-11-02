package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// Before clear all user data and create test user for preparation of test
func (t *APITest) Before() {
	admin, err := services.PasswordSignIn(services.Secret.Admin.UserName, services.Secret.Admin.Password)
	t.Assert(err == nil)
	services.Stop()
	services.Purge(admin)
	services.Start()
	services.PasswordSignUp("test@example.com", "test", "password", 10, entities.Normal)
}

// registerTestUser returns jwt containing created user. It's utility method for user action test
func (t *APITest) registerTestUser(id string, password string, hue int) (string, jwt.MapClaims) {
	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		Hue      int    `json:"hue"`
	}{
		ID:       id,
		Password: password,
		Hue:      hue,
	})
	t.Assert(err == nil)
	t.Post("/api/v1/register", "application/json", bytes.NewReader(obj))
	t.AssertOk()

	jwtObj := struct {
		Jwt string `json:"jwt"`
	}{}
	json.Unmarshal(t.ResponseBody, &jwtObj)

	token, err := jwt.Parse(jwtObj.Jwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(services.Secret.Auth.Salt), nil
	})
	t.Assert(err == nil)

	return jwtObj.Jwt, token.Claims.(jwt.MapClaims)
}

// TestRegister try sign in as new user
func (t *APITest) TestRegister() {
	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		Hue      int    `json:"hue"`
	}{
		ID:       "user@example.com",
		Password: "password",
		Hue:      120,
	})
	t.Assert(err == nil)
	t.Post("/api/v1/register", "application/json", bytes.NewReader(obj))
	t.AssertOk()
}

// TestRegisterInvalid try sign in as invalid user
func (t *APITest) TestRegisterInvalid() {
	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
		Hue      int    `json:"hue"`
	}{
		ID:       "user@example.com",
		Password: "password",
		Hue:      360,
	})
	t.Assert(err == nil)
	t.Post("/api/v1/register", "application/json", bytes.NewReader(obj))
	t.AssertStatus(422)
}

// TestLogin try login as admin
func (t *APITest) TestLogin() {
	t.registerTestUser("test@example.com", "password", 20)

	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{
		ID:       "test@example.com",
		Password: "password",
	})
	t.Assert(err == nil)
	t.Post("/api/v1/login", "application/json", bytes.NewReader(obj))
	t.AssertOk()
}

// TestLoginFailed try login as invalid user
func (t *APITest) TestLoginFailed() {
	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{})
	t.Assert(err == nil)
	t.Post("/api/v1/login", "application/json", bytes.NewReader(obj))
	t.AssertStatus(401)
}

// TestGetSettings try getting user private settings
func (t *APITest) TestGetSettings() {
	jwt, _ := t.registerTestUser("test@example.com", "password", 20)

	req, _ := http.NewRequest("GET", "/api/v1/settings", nil)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	t.NewTestRequest(req).Send()
	t.AssertOk()
}

// TestGetSettingsInvalid try getting user private settings without bearer token
func (t *APITest) TestGetSettingsInvalid() {
	t.registerTestUser("test@example.com", "password", 20)

	t.Get("/api/v1/settings")
	t.AssertStatus(401)
}

// TestChangeSettingsNoAuth try changing settings withotu authentication
func (t *APITest) TestChangeSettingsNoAuth() {
	t.registerTestUser("test@example.com", "password", 20)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{Value: "changed"})
	t.Assert(err == nil)
	t.Post("/api/v1/settings/cname", "application/json", bytes.NewBuffer(buf))
	t.AssertStatus(401)
}

// TestChangeSettingsCName try changing custom user name
func (t *APITest) TestChangeSettingsCName() {
	jwt, _ := t.registerTestUser("test@example.com", "password", 20)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{Value: "changed"})
	t.Assert(err == nil)
	req, _ := http.NewRequest("POST", "/api/v1/settings/cname", bytes.NewBuffer(buf))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	req.Header.Set("Content-Type", "application/json")
	t.NewTestRequest(req).Send()
	t.AssertOk()
}

// TestChangeSettingsCNameInvalid try changing invalid custom user name
func (t *APITest) TestChangeSettingsCNameInvalid() {
	jwt, _ := t.registerTestUser("test@example.com", "password", 20)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{})
	t.Assert(err == nil)
	req, _ := http.NewRequest("POST", "/api/v1/settings/cname", bytes.NewBuffer(buf))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	req.Header.Set("Content-Type", "application/json")
	t.NewTestRequest(req).Send()
	t.AssertStatus(422)
}

// TestChangeSettingsUseCName try changing using custom user name
func (t *APITest) TestChangeSettingsUseCName() {
	jwt, _ := t.registerTestUser("test@example.com", "password", 20)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{Value: "true"})
	t.Assert(err == nil)
	req, _ := http.NewRequest("POST", "/api/v1/settings/use_cname", bytes.NewBuffer(buf))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	req.Header.Set("Content-Type", "application/json")
	t.NewTestRequest(req).Send()
	t.AssertOk()
}

// TestChangeSettingsUseCNameInvalid try changing invalid using custom user name
func (t *APITest) TestChangeSettingsUseCNameInvalid() {
	jwt, _ := t.registerTestUser("test@example.com", "password", 20)
	buf, err := json.Marshal(struct {
		Value string `json:"value"`
	}{})
	t.Assert(err == nil)
	req, _ := http.NewRequest("POST", "/api/v1/settings/use_cname", bytes.NewBuffer(buf))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", jwt))
	req.Header.Set("Content-Type", "application/json")
	t.NewTestRequest(req).Send()
	t.AssertStatus(422)
}
