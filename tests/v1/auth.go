package v1

import (
	"bytes"
	"encoding/json"

	"github.com/yasshi2525/RushHour/app/services"
)

// TestLogin try login as admin
func (t *APITest) TestLogin() {
	obj, err := json.Marshal(struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{
		ID:       services.Secret.Admin.UserName,
		Password: services.Secret.Admin.Password,
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

// TestRegister try signin as new user
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

// TestRegisterInvalid try signin as invalid user
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
