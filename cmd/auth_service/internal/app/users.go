package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
)

func (a *App) RegisterUser(context *gin.Context) {
	user := types.NewUser()

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if user.Role == "" {
		user.Role = "popug"
	}

	err := a.repo.CreateUser(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := messages.UserMessage{
		Type:      messages.UserCreated,
		UserData:  user,
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res))
	//

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}

type UpdateRoleRequest struct {
	Email   string `json:"email"`
	NewRole string `json:"role"`
}

func (a *App) UpdateUserRole(context *gin.Context) {
	var req UpdateRoleRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user, err := a.repo.UpdateUserRole(req.Email, req.NewRole)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := messages.UserMessage{
		Type:      messages.UserRoleUpdated,
		UserData:  user,
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.beProducer.Send(string(res))
	//

	context.JSON(http.StatusOK, gin.H{"userId": user.ID, "email": user.Email, "role": user.Role})
}

func (a *App) GetAllPopugsIDs(context *gin.Context) {
	popugs, err := a.repo.GetAllPopugsIDs()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, popugs)
}
