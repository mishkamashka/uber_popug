package producer

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

type producer struct {
	cfg    *config
	writer sarama.AsyncProducer
}

func (p *producer) send(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: p.cfg.topic,
		Value: sarama.StringEncoder(message),
	}
	part, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", part)
	fmt.Println("Offset: ", o)
}

func (p *producer) OnStart(_ context.Context) error {
	if err := p.checkTopic(); err != nil {
		return err
	}

	var err error
	p.writer, err = p.newWriter()

	return err
}

func (p *producer) OnStop() {
	err := p.writer.Close()
	if err != nil {
		log.Println("close writer: ", err)
	}
}

func (p *producer) newWriter() (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Retry.Max = p.cfg.maxAttempts

	prod, err := sarama.NewAsyncProducer(p.cfg.brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("create producer: %w", err)
	}

	return prod, nil
}

func (p *producer) checkTopic() error {
	if p.cfg.topic == "" {
		return nil
	}

	var (
		topicIdx        int
		partitionsCount uint
		err             error
	)

	for i := 0; i < p.cfg.maxAttempts; i++ {
		if i != 0 {
			time.Sleep(time.Second)
		}

		partitionsCount, err = topicPartitionsCount(p.cfg.brokers, p.cfg.topic)
		if err != nil {
			continue
		}
		if partitionsCount == 0 {
			err = fmt.Errorf("no partitions for topic %s", p.cfg.topic)
			continue
		}
		topicIdx++

		break
	}

	return err
}

func topicPartitionsCount(brokers []string, topic string) (uint, error) {
	conn, err := sarama.NewClient(brokers, sarama.NewConfig())
	if err != nil {
		return 0, fmt.Errorf("dial %s: %w", brokers, err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Println("close kafka connection")
		}
	}()

	partitions, err := conn.Partitions(topic)
	if err != nil {
		return 0, fmt.Errorf("read partitions: %w", err)
	}

	return uint(len(partitions)), err
}
