package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
	"uber-popug/pkg/types/messages/v1"
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
	msg := v1.UserMessage{
		Type: v1.UserCreated,
		UserData: v1.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res), map[string]string{messages.V1: messages.V1})
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
	msg := v1.UserMessage{
		Type: v1.UserRoleUpdated,
		UserData: v1.UserData{
			ID:   user.ID,
			Role: user.Role,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.beProducer.Send(string(res), map[string]string{messages.Version: messages.V1})
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

type DeleteUserRequest struct {
	Email string `json:"email"`
}

func (a *App) DeleteUser(context *gin.Context) {
	var req DeleteUserRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if req.Email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty email"})
		context.Abort()
		return
	}

	user, err := a.repo.DeleteUser(req.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// send event
	msg := v1.UserMessage{
		Type: v1.UserDeleted,
		UserData: v1.UserData{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
		CreatedAt: time.Now(),
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.cudProducer.Send(string(res), map[string]string{messages.V1: messages.V1})
	//

	context.JSON(http.StatusOK, gin.H{"user_id": user.ID})
}

func (a *App) GetPopugEmail(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty userID"})
		context.Abort()
		return
	}

	email, err := a.repo.GetUsersEmail(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"email": email})
}
