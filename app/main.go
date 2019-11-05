package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/yasshi2525/RushHour/app/config"
	c1 "github.com/yasshi2525/RushHour/app/controllers/v1"
	"github.com/yasshi2525/RushHour/app/services"
)

func setupRouter(secret string) *gin.Engine {
	binding.Validator = new(c1.DefaultValidator)
	router := gin.New()
	store := cookie.NewStore([]byte(secret))
	router.Use(sessions.Sessions("mysession", store), gin.Logger(), gin.Recovery(), c1.GeneralHandler())

	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.LoadHTMLGlob("templates/*")

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/gamemap", c1.GameMap)
		apiV1.POST("/login", c1.Login)
		apiV1.POST("/register", c1.Register)
		apiV1.GET("/settings", c1.Settings)
		apiV1.POST("/settings/:resname", c1.ChangeSettings)
		apiV1.GET("/signout", c1.SignOut)
		apiV1.GET("/twitter", c1.Twitter)
		apiV1.GET("/google", c1.Google)
		apiV1.GET("/github", c1.GitHub)
		apiV1.GET("/twitterCallback", c1.TwitterCallback)
		apiV1.GET("/googleCallback", c1.GoogleCallback)
		apiV1.GET("/githubCallback", c1.GitHubCallback)
		apiV1.GET("/players", c1.Players)
		apiV1.POST("/rail_nodes", c1.Depart)
		apiV1.POST("/rail_nodes/extend", c1.Extend)
		apiV1.POST("/rail_nodes/connect", c1.Connect)
		apiV1.DELETE("/rail_nodes", c1.RemoveRailNode)
		apiV1.GET("/game", c1.GameStatus)
		apiV1.POST("/game/start", c1.StartGame)
		apiV1.POST("/game/stop", c1.StopGame)
		apiV1.DELETE("/game/purge", c1.PurgeUserData)
	}
	router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}

// @title RushHour REST API
// @version 1.0
// @description RushHour REST API
// @license.name MIT License
// @host rushhourgame.net
// @BasePath /api/v1
// @schemes https

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	conf, err := config.Load(fmt.Sprintf("%s/config", dir))
	if err != nil {
		panic(err)
	}
	services.Init(services.ServiceConfig{
		AppConf:   conf,
		IsPersist: true,
	})
	defer services.Terminate()

	services.Start()
	defer services.Stop()

	router := setupRouter(conf.Secret.Auth.Cookie)
	router.Run(":8080")
}
