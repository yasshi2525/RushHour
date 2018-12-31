package tests

import (
	"github.com/revel/revel/testing"
)

// AppTest is test suite for App
type AppTest struct {
	testing.TestSuite
}

// Before do nothing
func (t *AppTest) Before() {
	println("Set up")
}

// TestThatIndexPageWorks test
func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

// After do nothing
func (t *AppTest) After() {
	println("Tear down")
}
