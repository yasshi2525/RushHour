package tests

import (
	"github.com/revel/revel/testing"
)

// APITest is test suite for REST API
type APITest struct {
	testing.TestSuite
}

// Before starts game
func (t *APITest) Before() {
}

// TestThatIndexPageWorks test
func (t *APITest) TestThatIndexPageWorks() {
	t.Get("/api/v1/gamemap")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

// After stops game
func (t *APITest) After() {
}
