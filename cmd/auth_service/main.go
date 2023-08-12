package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"uber-popug/cmd/auth_service/internal/api"
	"uber-popug/cmd/auth_service/internal/middlewares"
	"uber-popug/cmd/auth_service/internal/repository"
)

func main() {
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	app := api.NewApp(repo)

	// Initialize Router
	router := initRouter(app)
	router.Run(":8080")
}
func initRouter(app *api.App) *gin.Engine {
	router := gin.Default()
	group := router.Group("/api")
	{
		group.POST("/token", app.GenerateToken)
		group.POST("/user/register", app.RegisterUser)
		secured := group.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", app.Ping)
		}
	}
	return router
}