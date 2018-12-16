package app

import (
	"github.com/yasshi2525/RushHour/app/models/entities"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// Db connection to the database
	Db *gorm.DB
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	revel.OnAppStart(InitDB, 1)
	revel.OnAppStart(MigrateDB, 2)
	revel.OnAppStart(InitGame, 3)
	// revel.OnAppStart(FillCache)

	revel.OnAppStop(CloseDB)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}

// InitDB connect database
func InitDB() {
	var err error

	driver, found := revel.Config.String("db.driver")

	if !found {
		panic("db.drvier is not defined")
	}

	spec, found := revel.Config.String("db.spec")

	if !found {
		panic("db.spec is not defined")
	}

	Db, err = gorm.Open(driver, spec)
	Db.LogMode(true)

	if err != nil {
		panic("failed to connect database")
	}

	revel.RevelLog.Info("connect database successfully")
}

// MigrateDB migrate database
func MigrateDB() {
	Db.AutoMigrate(
		&entities.Company{},
		&entities.Residence{},
		&entities.Human{},
		&entities.Player{},
		&entities.RailNode{},
		&entities.RailEdge{},
		&entities.Platform{},
		&entities.Gate{},
		&entities.Station{},
	)
}

// CloseDB close database connection
func CloseDB() {
	if err := Db.Close(); err != nil {
		revel.AppLog.Error("failed to close the database", "error", err)
	}
	revel.RevelLog.Info("disconnect database successfully")
}

// InitGame setup RushHour envirionment
func InitGame() {
	var admin entities.Player
	var residence entities.Residence
	var company entities.Company
	var human entities.Human
	var railNode entities.RailNode

	Db.Where(entities.Player{DisplayName: "Admin", Password: "encodedPassword"}).FirstOrCreate(&admin)
	revel.AppLog.Info("created Admin player")

	Db.FirstOrCreate(
		&residence,
		entities.Residence{
			Point: entities.Point{X: 5, Y: 5},
		})
	Db.FirstOrCreate(
		&company,
		entities.Company{
			Point: entities.Point{X: 10, Y: 10},
		})
	Db.FirstOrCreate(
		&human,
		entities.Human{
			Point:     entities.Point{X: 10, Y: 10},
			FromRefer: residence.ID,
			ToRefer:   company.ID,
			On:        entities.OnGround,
		})

	revel.AppLog.Info("created Public Facilities")

	Db.FirstOrCreate(
		&railNode,
		entities.RailNode{
			Owner: entities.Owner{OwnerRefer: admin.ID},
			Point: entities.Point{X: 100, Y: 100},
		})

	revel.AppLog.Info("created Private Facilities")
}
