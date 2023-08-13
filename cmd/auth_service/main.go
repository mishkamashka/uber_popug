package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"uber-popug/cmd/auth_service/internal/app"
	"uber-popug/cmd/auth_service/internal/repository"
	"uber-popug/pkg/kafka/producer"
	middlewares "uber-popug/pkg/middlewares"
)

var (
	brokers = []string{"localhost:9092"}
)

func main() {
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	cudProducerConfig := producer.NewConfig(brokers, "users-stream", "auth-service")
	beProducerConfig := producer.NewConfig(brokers, "users", "auth-service")

	cudProducer, err := producer.NewProducer(cudProducerConfig)
	if err != nil {
		log.Fatal(err)
	}

	beProducer, err := producer.NewProducer(beProducerConfig)
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(repo, cudProducer, beProducer)

	// Initialize Router
	router := initRouter(app)
	router.Run(":2400")
}
func initRouter(app *app.App) *gin.Engine {
	router := gin.Default()
	group := router.Group("/api")
	{
		group.POST("/token", app.GenerateToken)
		group.POST("/user/register", app.RegisterUser)
		secured := group.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", app.Ping)
		}

		admin := group.Group("/admin").Use(middlewares.AdminAuth())
		{
			admin.POST("/user/role", app.UpdateUserRole)
			admin.DELETE("/user", app.DeleteUser)
		}

		internal := group.Group("/internal").Use()
		{
			internal.GET("/popugs", app.GetAllPopugsIDs)
		}
	}
	return router
}
