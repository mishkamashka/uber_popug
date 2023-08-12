package app

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"uber-popug/cmd/auth_service/internal/types"
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
	// user created
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
	// role updated
	//

	context.JSON(http.StatusOK, gin.H{"userId": user.ID, "email": user.Email, "role": user.Role})
}
