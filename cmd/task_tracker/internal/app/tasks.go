package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"log"
	"net/http"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorId   string `json:"creator_id"`
}

func (a *App) CreateTask(context *gin.Context) {
	req := CreateTaskRequest{}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// get random popug

	id, _ := uuid.GenerateUUID()

	task := &types.Task{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Status:      "open",
		AssigneeId:  "",
		CreatorId:   req.CreatorId,
	}

	task.GeneratePrices()

	err := a.repo.CreateTask(task)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := messages.TaskMessage{
		Type:     messages.UserCreated,
		TaskData: task,
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res))
	//

	context.JSON(http.StatusCreated, gin.H{
		"taskId":            task.ID,
		"assignee_id":       task.AssigneeId,
		"price_for_assign":  task.PriceForAssign,
		"price_for_closing": task.PriceForClosing,
	})
}

type CloseTaskRequest struct {
	TaskID string `json:"task_id"`
}

func (a *App) UpdateUserRole(context *gin.Context) {
	var req CloseTaskRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user, err := a.repo.UpdateTaskStatus(req.TaskID, "closed")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := messages.TaskMessage{
		Type:     messages.UserRoleUpdated,
		TaskData: task,
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.beProducer.Send(string(res))
	//

	context.JSON(http.StatusOK, gin.H{"task_id": task.ID, "status": task.Status})
}
