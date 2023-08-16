package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"log"
	"math/rand"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
	"uber-popug/pkg/types/messages/v1"
)

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *App) CreateTask(context *gin.Context) {
	userID, ok := context.Value("userID").(string)
	if userID == "" || !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no userID or userID is not string"})
		context.Abort()
		return
	}

	req := CreateTaskRequest{}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// get all popugs (seems like too much time spent here but no time to think for a better solution)
	popugs, err := a.client.GetAllPopugsIDs()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	id, _ := uuid.GenerateUUID()

	task := &types.Task{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Status:      "open",
		AssigneeId:  popugs[rand.Intn(len(popugs))],
		CreatorId:   userID,
	}

	task.GeneratePrices()

	err = a.repo.CreateTask(task)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := v1.TaskMessage{
		Type: v1.TaskCreated,
		Data: v1.TaskData{
			ID:              task.ID,
			Name:            task.Name,
			Description:     task.Description,
			Status:          task.Status,
			PriceForAssign:  task.PriceForAssign,
			PriceForClosing: task.PriceForClosing,
			AssigneeId:      task.AssigneeId,
			CreatorId:       task.CreatorId,
			CreatedAt:       task.CreatedAt,
			UpdatedAt:       task.UpdatedAt,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res), map[string]string{messages.Version: messages.V1})
	//

	context.JSON(http.StatusCreated, gin.H{
		"taskId":            task.ID,
		"assignee_id":       task.AssigneeId,
		"price_for_assign":  task.PriceForAssign,
		"price_for_closing": task.PriceForClosing,
	})
}
