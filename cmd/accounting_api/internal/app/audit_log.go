package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/util"
)

func (a *App) CreateAuditLog(task *types.Task) error {

	return a.repo.CreateAuditLog(audilLog)
}

func (a *App) GetPopugTodayAuditLog(context *gin.Context) {
	userID, ok := context.Value("userID").(string)
	if userID == "" || !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no userID or userID is not string"})
		context.Abort()
		return
	}

	to := time.Now()

	from := util.TruncateToDay(to)

	balance, err := a.repo.GetUserAuditLogsForPeriod(userID, from, to)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, balance)
}
