package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"uber-popug/pkg/types/messages"
	v2 "uber-popug/pkg/types/messages/v2"
)

type app interface {
	ReassignUsersTasks(userID string) error
}

type handler struct {
	app app
}

func NewHandler(app app) *handler {
	return &handler{
		app: app,
	}
}

func (h *handler) Handle(msg *sarama.ConsumerMessage) error {
	version := ""
	for _, v := range msg.Headers {
		if string(v.Key) == messages.Version {
			version = string(v.Value)
		}
	}

	switch version {
	case "", messages.V1:
		return h.handleV2Msg(msg)
	default:
		log.Printf("unknown message version: %s", version)
	}

	return nil
}

func (h *handler) handleV2Msg(msg *sarama.ConsumerMessage) error {
	var event v2.TaskMessage

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return fmt.Errorf("unmarshalling msg: %s", err)
	}

	switch event.Type {
	case v2.TaskCreated:
		// minus amount
	case v2.TaskClosed:
		// add amount

	}
	return nil
}
