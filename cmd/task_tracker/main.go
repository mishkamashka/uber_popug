package main

import (
	"log"
	"uber-popug/cmd/task_tracker/internal/api"
	"uber-popug/cmd/task_tracker/internal/app"
)

var (
	brokers = []string{"localhost:9092"}
)

func main() {
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	//cudProducerConfig := producer.NewConfig(brokers, "tasks-stream", "tasks-service")
	//beProducerConfig := producer.NewConfig(brokers, "tasks", "tasks-service")
	//
	//cudConsumerConfig := consumer.NewConfig(brokers, []string{"users-stream"}, "tasks-service")
	//beConsumerConfig := consumer.NewConfig(brokers, []string{"users"}, "tasks-service")
	//
	//c, err := consumer.New(cudConsumerConfig)
	//c.OnMessage()

	//cudProducer, err := producer.NewProducer(cudProducerConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//beProducer, err := producer.NewProducer(beProducerConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	app := app.NewApp(repo, cudProducer, beProducer)

	// Initialize Router
	router := api.NewApi(app)
	router.Run(":8080")
}
