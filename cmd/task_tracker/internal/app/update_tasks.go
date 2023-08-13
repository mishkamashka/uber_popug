package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
	"uber-popug/pkg/types/messages"
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

	task, err := a.repo.UpdateTaskStatus(req.TaskID, "closed")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := messages.TaskMessage{
		Type:      messages.TaskClosed,
		Data:      task,
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.beProducer.Send(string(res))
	//

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

		msg := messages.TaskMessage{
			Type:      messages.TaskReassigned,
			Data:      task,
			CreatedAt: time.Now(),
		}
		res, err := json.Marshal(msg)
		if err != nil {
			log.Println("error producing message")
		}

		a.cudProducer.Send(string(res))
	}
}

func (a *App) ReassignUsersTasks(userID string) error {
	tasks, err := a.repo.GetUserTasks(userID)
	if err != nil {
		return fmt.Errorf("get user's tasks: %s", err)
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

		msg := messages.TaskMessage{
			Type:      messages.TaskReassigned,
			Data:      task,
			CreatedAt: time.Now(),
		}
		res, err := json.Marshal(msg)
		if err != nil {
			log.Println("error producing message")
		}

		a.cudProducer.Send(string(res))
	}

	return nil
}
