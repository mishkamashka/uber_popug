package main

import (
	"context"
	"log"
	"uber-popug/cmd/accounting/internal/accounting_handler"
	"uber-popug/cmd/accounting/internal/api"
	"uber-popug/cmd/accounting/internal/app"
	"uber-popug/cmd/accounting/internal/repository"
	"uber-popug/pkg/kafka/consumer"
)

var (
	brokers = []string{"localhost:9092"}
)

func main() {
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(repo)

	// tasks' events consumer
	cudConsumerConfig := consumer.NewConfig(brokers, []string{"tasks-stream"}, "accounting-service")

	c, err := consumer.New(cudConsumerConfig)
	c.OnMessage(accounting_handler.NewHandler(app).Handle)
	c.OnStart(context.Background())

	// Initialize Router
	router := api.NewApi(app)
	router.Run(":2402")
}
