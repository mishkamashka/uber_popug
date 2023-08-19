package accounting_handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
	v1 "uber-popug/pkg/types/messages/v1"
	v2 "uber-popug/pkg/types/messages/v2"
)

type app interface {
	CreateTaskClosedAuditLog(task *types.Task) error
	CreateTaskAssignedAuditLog(task *types.Task) error

	UpdatePopugBalance(userID string, amount int) error
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
	case "", messages.V2:
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

	task := taskDataV2ToTask(event.Data)

	switch event.Type {
	case v2.TaskCreated, v2.TaskReassigned:
		err := h.app.CreateTaskAssignedAuditLog(task)
		if err != nil {
			log.Println("create audit log: " + err.Error())
		}

		err = h.app.UpdatePopugBalance(task.AssigneeId, -int(task.PriceForAssign))
		if err != nil {
			log.Println("update popug's audit log: " + err.Error())
		}
	case v2.TaskClosed:
		err := h.app.CreateTaskClosedAuditLog(task)
		if err != nil {
			log.Println("create audit log: " + err.Error())
		}

		err = h.app.UpdatePopugBalance(task.AssigneeId, int(task.PriceForClosing))
		if err != nil {
			log.Println("update popug's audit log: " + err.Error())
		}
	}

	return nil
}

func taskDataV2ToTask(data v2.TaskData) *types.Task {
	return &types.Task{
		ID:              data.ID,
		Title:           data.Title,
		JiraID:          data.JiraID,
		Description:     data.Description,
		Status:          data.Status,
		PriceForAssign:  data.PriceForAssign,
		PriceForClosing: data.PriceForClosing,
		AssigneeId:      data.AssigneeId,
		CreatorId:       data.CreatorId,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}

func (h *handler) handleV1Msg(msg *sarama.ConsumerMessage) error {
	var event v1.TaskMessage

	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return fmt.Errorf("unmarshalling msg: %s", err)
	}

	task := taskDataV1ToTask(event.Data)

	switch event.Type {
	case v1.TaskCreated, v1.TaskReassigned:
		err := h.app.CreateTaskAssignedAuditLog(task)
		if err != nil {
			log.Println("create audit log: " + err.Error())
		}

		err = h.app.UpdatePopugBalance(task.AssigneeId, -int(task.PriceForAssign))
		if err != nil {
			log.Println("update popug's audit log: " + err.Error())
		}
	case v1.TaskClosed:
		err := h.app.CreateTaskClosedAuditLog(task)
		if err != nil {
			log.Println("create audit log: " + err.Error())
		}

		err = h.app.UpdatePopugBalance(task.AssigneeId, int(task.PriceForClosing))
		if err != nil {
			log.Println("update popug's audit log: " + err.Error())
		}
	}

	return nil
}

func taskDataV1ToTask(data v1.TaskData) *types.Task {
	return &types.Task{
		ID:              data.ID,
		Title:           data.Name,
		Description:     data.Description,
		Status:          data.Status,
		PriceForAssign:  data.PriceForAssign,
		PriceForClosing: data.PriceForClosing,
		AssigneeId:      data.AssigneeId,
		CreatorId:       data.CreatorId,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
	}
}
