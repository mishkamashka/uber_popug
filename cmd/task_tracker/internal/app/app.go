package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateTask(user *types.Task) error
	GetUserTasks(userID string) ([]*types.Task, error)
	UpdateTaskStatus(taskID, status string) (*types.Task, error)
	GetAllOpenTasks() ([]*types.Task, error)
	UpdateTask(task *types.Task) error
}

type producer interface {
	Send(msg string)
}

type usersClient interface {
	GetAllPopugsIDs() ([]string, error)
}
type App struct {
	repo        repository
	client      usersClient
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
