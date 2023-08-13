package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"uber-popug/pkg/types/messages"
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
	var event messages.UserMessage

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return fmt.Errorf("unmarshalling msg: %s", err)
	}

	if event.Type == messages.UserDeleted {
		err := h.app.ReassignUsersTasks(event.UserData.ID)
		if err != nil {
			log.Printf("reassigning user's tasks: %s", err)
		}

		log.Println("reassigned deleted user's tasks")
	}

	return nil
}
