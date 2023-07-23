package consumer

import (
	"github.com/IBM/sarama"
	"log"
)

type handler struct {
	handleFunc handleFn
	ready      chan bool
}

func NewHandler() *handler {
	return &handler{
		ready: make(chan bool),
	}
}

func (h *handler) WithHandleFunc(fn handleFn) {
	h.handleFunc = fn
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *handler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the handler as ready
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a handler loop of handlerGroupClaim's Messages().
func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine.
	for {
		select {
		case message := <-claim.Messages():
			err := h.handleFunc(message)
			log.Printf("processing message: %s", err)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
