package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uber-popug/pkg/types"
)

type repository interface {
	CreateAuditLog(log *types.AuditLog) error
	GetUserAuditLogsForPeriod(userID string, from, to time.Time) ([]*types.AuditLog, error)

	GetPopugBalance(userID string) (*types.Balance, error)
	UpdatePopugBalanceByValue(userID string, amount int) error
	GetAllNegativePopugsBalances() ([]*types.Balance, error)
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
