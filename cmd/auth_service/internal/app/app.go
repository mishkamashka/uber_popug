package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateUser(user *types.User) error
	GetUserByEmail(email string) (*types.User, error)
	UpdateUserRole(email, role string) (*types.User, error)
}

type producer interface {
	Send(msg string)
}

type App struct {
	repo        repository
	cudProducer producer
	beProducer  producer
}

func NewApp(repo repository, cudProducer, beProducer producer) *App {
	return &App{
		repo:        repo,
		cudProducer: cudProducer,
		beProducer:  beProducer,
	}
}

func (a *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
