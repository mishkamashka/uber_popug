package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"uber-popug/pkg/types/messages"
	"uber-popug/pkg/types/messages/v1"
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
		return h.handleV1Msg(msg)
	default:
		log.Printf("unknown message version: %s", version)
	}

	return nil
}

func (h *handler) handleV1Msg(msg *sarama.ConsumerMessage) error {
	var event v1.UserMessage

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return fmt.Errorf("unmarshalling msg: %s", err)
	}

	if event.Type == v1.UserDeleted ||
		(event.Type == v1.UserRoleUpdated && event.UserData.Role != "popug") {
		err := h.app.ReassignUsersTasks(event.UserData.ID)
		if err != nil {
			log.Printf("reassigning user's tasks: %s", err)
		}

		log.Println("reassigned deleted user's tasks")
	}

	return nil
}
