package main

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"

	"github.com/yasshi2525/RushHour/auth"
	"github.com/yasshi2525/RushHour/config"
	"github.com/yasshi2525/RushHour/controllers"
	v1 "github.com/yasshi2525/RushHour/controllers/v1"
	"github.com/yasshi2525/RushHour/services"
)

// @title RushHour REST API
// @version 1.0
// @description RushHour REST API
// @license.name MIT License
// @host rushhourgame.net
// @BasePath /api/v1
// @schemes https

var readiness string

func loadConf() (*config.Config, error) {
	if dir, err := os.Getwd(); err != nil {
		return nil, err
	} else if conf, err := config.Load(fmt.Sprintf("%s/config", dir)); err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}

func setupRouter(secret string) *gin.Engine {
	binding.Validator = new(v1.DefaultValidator)
	router := gin.New()
	router.Use(gin.Recovery())

	// healthCheck
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/readiness", func(c *gin.Context) {
		if readiness != "" {
			c.String(http.StatusServiceUnavailable, readiness)
		} else {
			c.String(http.StatusOK, "OK")
		}
	})

	store := cookie.NewStore([]byte(secret))
	app := router.Group("/", gin.Logger(), sessions.Sessions("rushhour", store))

	app.Static("/assets", "./assets")
	app.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.LoadHTMLGlob("templates/*")

	// index always return html
	index := app.Group("/")
	{
		index.GET("/", controllers.Index)
		index.POST("/", controllers.Index)
	}

	// redirecting page for OAuth
	// it might causes err in invalid configuration
	oauth := app.Group("/", controllers.OAuthHandler())
	{
		oauth.GET("/twitter", controllers.Twitter)
		oauth.GET("/google", controllers.Google)
		oauth.GET("/github", controllers.GitHub)
	}

	// callback page from OAuth
	// it might causes err in invalid configuration
	callback := app.Group("/", controllers.CallbackHandler(), controllers.RegisterHandler())
	{
		callback.GET("/google/callback", controllers.GoogleCallback)
		callback.GET("/github/callback", controllers.GitHubCallback)
	}
	// twitter callback is irregular pattern
	app.GET("/twitter/callback", controllers.TwitterCallback, controllers.RegisterHandler())

	api := app.Group("/api/v1")
	{
		// available only in operation
		ops := api.Group("/", v1.MaintenanceHandler())
		{
			// no need auth (only under operation)
			shared := ops.Group("/", v1.ModelHandler())
			{
				shared.GET("/gamemap", v1.GameMap)
				shared.GET("/players", v1.Players)
				shared.POST("/register", v1.Register)
			}

			// need user authorization (only under operation)
			user := ops.Group("/", v1.JWTHandler(), v1.ModelHandler())
			{
				user.GET("/settings", v1.Settings)
				user.POST("/settings/:resname", v1.ChangeSettings)
				user.POST("/signout", v1.SignOut)
				user.POST("/rail_nodes", v1.Depart)
				user.POST("/rail_nodes/extend", v1.Extend)
				user.POST("/rail_nodes/connect", v1.Connect)
				user.DELETE("/rail_nodes", v1.RemoveRailNode)
			}
		}

		// always available even though under maintenance
		always := api.Group("/")
		{
			// no need auth (always)
			shared := always.Group("/", v1.ModelHandler())
			{
				shared.POST("/login", v1.Login) // forbit normal user under maintenance
				shared.GET("/game", v1.GameStatus)
			}
			// need administrator authorization (always)
			admin := always.Group("/", v1.JWTHandler(), v1.AdminHandler(), v1.ModelHandler())
			{
				admin.POST("/game/start", v1.StartGame)
				admin.POST("/game/stop", v1.StopGame)
				admin.DELETE("/game/purge", v1.PurgeUserData)
			}
		}
	}
	return router
}

func main() {
	// randomization
	if seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64)); err != nil {
		panic(err)
	} else {
		rand.Seed(seed.Int64())
	}
	// configuration
	if conf, err := loadConf(); err != nil {
		panic(err)
	} else if auther, err := auth.GetAuther(conf.Secret.Auth); err != nil {
		panic(err)
	} else {
		// run server
		router := setupRouter(conf.Secret.Auth.Cookie)
		router.Run(":8080")

		readiness = "initializing ..."

		// prepare service
		services.Init(conf, auther)
		defer services.Terminate()

		readiness = "starting ..."

		services.Start()
		defer services.Stop()

		// prepare controller
		controllers.InitController(auther)
		v1.InitController(conf, auther)

		readiness = ""
		defer func() {
			readiness = "shut down ..."
		}()
	}
}
