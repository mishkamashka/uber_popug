package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/cmd/auth_service/internal/types"
)

type repository interface {
	CreateUser(user *types.User) error
	GetUserByEmail(email string) (*types.User, error)
	UpdateUserRole(email, role string) (*types.User, error)
}

type App struct {
	repo repository
}

func NewApp(repo repository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
