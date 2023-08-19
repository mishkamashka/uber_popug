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

type mailSender interface {
	Send(email string, msg []byte)
}

type app struct {
	accountingClient accountingClient
	tasksClient      tasksClient
	mailSender       mailSender
}

func NewApp(accClient accountingClient, tasksClient tasksClient, sender mailSender) *app {
	return &app{
		accountingClient: accClient,
		tasksClient:      tasksClient,
		mailSender:       sender,
	}
}

func (a *app) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
