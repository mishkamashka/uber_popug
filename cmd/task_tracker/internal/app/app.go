package app

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
	"uber-popug/cmd/task_tracker/internal/popug_client"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateTask(user *types.Task) error
	GetUserTasks(userID string) ([]*types.Task, error)
	CloseTask(taskID string) (*types.Task, error)
	GetAllOpenTasks() ([]*types.Task, error)
	UpdateTask(task *types.Task) error
	DeleteTask(taskID string) error
	TopTask(from time.Time) (*types.Task, error)
	GetAssignedTasksFromTime(from time.Time) ([]*types.Task, error)
	GetClosedTasksFromTime(from time.Time) ([]*types.Task, error)
	GetActiveTasksFromPeriod(from, to time.Time) ([]*types.Task, error)
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
	rand        *rand.Rand
}

func NewApp(repo repository, cudProducer, beProducer producer) *App {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return &App{
		repo:        repo,
		cudProducer: cudProducer,
		beProducer:  beProducer,
		client:      popug_client.New(),
		rand:        r,
	}
}

func (a *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
