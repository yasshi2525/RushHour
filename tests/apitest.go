package tests

import (
	"github.com/revel/revel/testing"
	"github.com/yasshi2525/RushHour/app"
	"github.com/yasshi2525/RushHour/app/services"
)

// ApiTest is test suite for REST API
type ApiTest struct {
	testing.TestSuite
}

// Before starts game
func (t *ApiTest) Before() {
	println("Set up")
	services.LoadConf()
	services.InitStorage()
	services.InitPersistence()
	app.InitGame()
}

// TestThatIndexPageWorks test
func (t *ApiTest) TestThatIndexPageWorks() {
	t.Get("/api/v1/gamemap")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

// After stops game
func (t *ApiTest) After() {
	println("Tear down")

	app.StopGame()
	services.Backup()
	services.TerminatePersistence()
}
