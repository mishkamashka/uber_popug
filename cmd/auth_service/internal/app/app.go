package app

import (
	"context"
	"net/http"
	"uber-popug/cmd/auth_service/internal/jwt"

	"github.com/gin-gonic/gin"

	"uber-popug/cmd/auth_service/internal/types"
)

type repository interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
}

type App struct {
	repo repository
}

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

	err := a.repo.CreateUser(context, user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) GenerateToken(context *gin.Context) {
	var request TokenRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	user, err := a.repo.GetUserByEmail(context, request.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}
	tokenString, err := jwt.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
