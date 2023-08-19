package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"uber-popug/pkg/types/messages"
	v1 "uber-popug/pkg/types/messages/v1"
)

type DeleteTaskRequest struct {
	TaskID string `json:"task_id"`
}

func (a *App) DeleteTask(context *gin.Context) {
	var req DeleteTaskRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	err := a.repo.DeleteTask(req.TaskID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := v1.TaskMessage{
		Type: v1.TaskDeleted,
		Data: v1.TaskData{
			ID: req.TaskID,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res), map[string]string{messages.Version: messages.V1})
	//

	context.JSON(http.StatusOK, gin.H{"taskID": req.TaskID})
}
