package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uber-popug/pkg/types"
)

type accountingClient interface {
	Checkout(userID string, dayTotal int) error
}

type tasksClient interface {
	GetAllUpdatedTasksForToday() ([]*types.Task, error)
}

type userClient interface {
	GetUserEmail(userID string) (string, error)
}

type mailSender interface {
	Send(email string, dayTotal int) error
}

type app struct {
	accountingClient accountingClient
	tasksClient      tasksClient
	usersClient      userClient
	mailSender       mailSender
}

func NewApp(accClient accountingClient, tasksClient tasksClient, usersClient userClient, sender mailSender) *app {
	return &app{
		accountingClient: accClient,
		tasksClient:      tasksClient,
		usersClient:      usersClient,
		mailSender:       sender,
	}
}

func (a *app) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
