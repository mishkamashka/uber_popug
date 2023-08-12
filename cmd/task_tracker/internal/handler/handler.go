package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"uber-popug/pkg/types/messages"
)

type app interface {
	ReassignUserTasks(userID string) error
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
		err := h.app.ReassignUserTasks(event.UserData.ID)
		if err != nil {
			log.Printf("deleting user's tasks: %s", err)
		}
	}

	return nil
}
