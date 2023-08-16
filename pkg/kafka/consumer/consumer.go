package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type consumer struct {
	handler *handler
	cfg     *config
	ready   chan bool
	done    chan struct{}
}

func New(cfg *config) (*consumer, error) {
	c := consumer{
		handler: NewHandler(),
		cfg:     cfg,
		done:    make(chan struct{}),
		ready:   make(chan bool),
	}

	return &c, nil
}

func (c *consumer) OnMessage(fn handleFn) {
	if fn != nil {
		c.handler.WithHandleFunc(func(msg *sarama.ConsumerMessage) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("consumer handler: %s", r)
				}
			}()
			return fn(msg)
		})
	}
}

func (c *consumer) OnStart(ctx context.Context) error {
	if c.handler.handleFunc == nil {
		return errors.New("handleFn is nil")
	}

	go c.run(ctx)
	return nil
}

func (c *consumer) OnStop() {}

func (c *consumer) run(ctx context.Context) {
	for i := 1; true; i++ {
		err := c.consumeFromGroup(ctx)
		if err != nil && i == 3 {
			log.Printf("consume from group: try %d", i)
			i = 0
		}
		select {
		case <-ctx.Done():
			close(c.done)
			return
		case <-time.After(5 * time.Second):
		}
	}
}

func (c *consumer) consumeFromGroup(ctx context.Context) error {
	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	config := sarama.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(c.cfg.brokers, c.cfg.groupID, config)
	if err != nil {
		cancel()
		return fmt.Errorf("create consumer group: %w", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, c.cfg.topics, c.handler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	return nil
}
