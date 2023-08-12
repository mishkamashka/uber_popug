package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateTask(user *types.Task) error
	GetUserTasks(email string) (*types.Task, error)
	ReassignTasks() error
	UpdateTaskStatus(taskID, status string) (*types.Task, error)
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
