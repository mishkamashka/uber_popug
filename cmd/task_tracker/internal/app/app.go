package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/cmd/task_tracker/internal/popug_client"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateTask(user *types.Task) error
	GetUserTasks(userID string) ([]*types.Task, error)
	UpdateTaskStatus(taskID, status string) (*types.Task, error)
	GetAllOpenTasks() ([]*types.Task, error)
	UpdateTask(task *types.Task) error
	DeleteTask(taskID string) error
}

type producer interface {
	Send(msg string, headers map[string]string)
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
		client:      popug_client.New(),
	}
}

func (a *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
