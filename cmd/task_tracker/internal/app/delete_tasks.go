package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	context.JSON(http.StatusOK, gin.H{"taskID": req.TaskID})
}
