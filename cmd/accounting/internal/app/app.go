package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"uber-popug/pkg/types"
)

type repository interface {
	CreateAuditLog(log *types.AuditLog) error
	GetUserAuditLogsForPeriod(userID string, from, to time.Time) ([]*types.AuditLog, error)

	GetPopugBalance(userID string) (*types.Balance, error)
	UpdatePopugBalanceByValue(userID string, amount int) error
	GetAllNegativePopugsBalances() ([]*types.Balance, error)
}

type producer interface {
	Send(msg string, headers map[string]string)
}

type App struct {
	repo       repository
	beProducer producer
}

func NewApp(repo repository, producer producer) *App {
	return &App{
		repo:       repo,
		beProducer: producer,
	}
}

func (a *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
