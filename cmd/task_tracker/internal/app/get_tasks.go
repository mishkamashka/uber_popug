package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uber-popug/pkg/util"
)

func (a *App) GetUserTasks(context *gin.Context) {
	userID, ok := context.Value("userID").(string)
	if userID == "" || !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no userID or userID is not string"})
		context.Abort()
		return
	}

	tasks, err := a.repo.GetUserTasks(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, tasks)
}

func (a *App) YesterdayTasks(context *gin.Context) {
	to := util.TruncateToDay(time.Now())
	from := to.AddDate(0, 0, -1)

	tasks, err := a.repo.GetActiveTasksFromPeriod(from, to)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
