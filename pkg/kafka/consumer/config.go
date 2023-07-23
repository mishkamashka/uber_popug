package consumer

import (
	"github.com/IBM/sarama"
	"time"
)

type handleFn = func(msg *sarama.ConsumerMessage) error

type config struct {
	brokers []string
	topics  []string
	groupID string

	// Default is 100
	bufferSize int
	// Default is 10e3
	minBytes int
	// Default is 10e6
	maxBytes int
	// Default is 3
	maxAttempts int
	// Default is 10 seconds
	readTimeout time.Duration
	// Default is 5 seconds. Used as an approximate value
	commitInterval time.Duration

	startOffset int64
}

func NewConfig(brokers, topics []string, groupID string) *config {
	return &config{
		brokers:        brokers,
		topics:         topics,
		groupID:        groupID,
		minBytes:       10e3,
		maxBytes:       10e6,
		maxAttempts:    3,
		readTimeout:    10 * time.Second,
		commitInterval: 5 * time.Second,
	}
}
