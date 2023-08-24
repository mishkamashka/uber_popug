package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"log"
	"math/rand"
	"net/http"
	"time"
	"uber-popug/pkg/types/messages"
	"uber-popug/pkg/types/messages/v1"
	v2 "uber-popug/pkg/types/messages/v2"
)

type CloseTaskRequest struct {
	TaskID string `json:"task_id"`
}

func (a *App) CloseTask(context *gin.Context) {
	var req CloseTaskRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	task, err := a.repo.CloseTask(req.TaskID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	id, _ := uuid.GenerateUUID()

	msg := v2.TaskMessage{
		ID:   id,
		Type: v2.TaskClosed,
		Data: v2.TaskData{
			ID:              task.ID,
			Title:           task.Title,
			JiraID:          task.JiraID,
			Status:          task.Status,
			PriceForClosing: task.PriceForClosing,
			AssigneeId:      task.AssigneeId,
			UpdatedAt:       task.UpdatedAt,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}

	a.cudProducer.Send(string(res), map[string]string{messages.Version: messages.V2})

	context.JSON(http.StatusOK, gin.H{"task_id": task.ID, "status": task.Status})
}

func (a *App) ReassignTasks(context *gin.Context) {
	tasks, err := a.repo.GetAllOpenTasks()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	popugs, err := a.client.GetAllPopugsIDs()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	for _, task := range tasks {
		task.AssigneeId = popugs[rand.Intn(len(popugs))]

		err := a.repo.UpdateTask(task)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		msg := v1.TaskMessage{
			Type: v1.TaskReassigned,
			Data: v1.TaskData{
				ID:              task.ID,
				Name:            task.Title,
				Description:     task.Description,
				Status:          task.Status,
				PriceForAssign:  task.PriceForAssign,
				PriceForClosing: task.PriceForClosing,
				AssigneeId:      task.AssigneeId,
				UpdatedAt:       task.UpdatedAt,
			},
			CreatedAt: time.Now(),
		}
		res, err := json.Marshal(msg)
		if err != nil {
			log.Println("error producing message")
		}

		a.cudProducer.Send(string(res), map[string]string{messages.Version: messages.V1})
	}
}

func (a *App) ReassignUsersTasks(userID string) error {
	tasks, err := a.repo.GetUserTasks(userID)
	if err != nil {
		return fmt.Errorf("get user's tasks: %s", err)
	}

	if len(tasks) == 0 {
		return nil
	}

	popugs, err := a.client.GetAllPopugsIDs()
	if err != nil {
		return fmt.Errorf("get popugs: %s", err)
	}

	for _, task := range tasks {
		task.AssigneeId = popugs[rand.Intn(len(popugs))]

		err := a.repo.UpdateTask(task)
		if err != nil {
			log.Printf("update task: %s", err)
			continue
		}

		msg := v1.TaskMessage{
			Type: v1.TaskReassigned,
			Data: v1.TaskData{
				ID:              task.ID,
				Name:            task.Title,
				Description:     task.Description,
				Status:          task.Status,
				PriceForAssign:  task.PriceForAssign,
				PriceForClosing: task.PriceForClosing,
				AssigneeId:      task.AssigneeId,
				UpdatedAt:       task.UpdatedAt,
			},
			CreatedAt: time.Now(),
		}
		res, err := json.Marshal(msg)
		if err != nil {
			log.Println("error producing message")
		}

		a.cudProducer.Send(string(res), map[string]string{messages.Version: messages.V1})
	}

	return nil
}
