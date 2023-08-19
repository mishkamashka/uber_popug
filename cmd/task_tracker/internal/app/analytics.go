package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uber-popug/pkg/util"
)

const (
	day   = "1d"
	week  = "1w"
	month = "1m"
)

func (a *App) TopTask(context *gin.Context) {
	period := context.Query("period")
	if period == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no period in query"})
		context.Abort()
		return
	}

	today := util.TruncateToDay(time.Now())

	var from time.Time

	switch period {
	case day:
		from = today
	case week:
		from = today.AddDate(0, 0, -6)
	case month:
		from = today.AddDate(0, -1, 0)
	default:
		context.JSON(http.StatusBadRequest, gin.H{"error": "unknown period" + period})
		context.Abort()
		return
	}

	task, err := a.repo.TopTask(from)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, task)
}

func (a *App) TodayEarnings(context *gin.Context) {
	today := util.TruncateToDay(time.Now())

	assignedTasks, err := a.repo.GetAssignedTasksFromTime(today)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	closedTasks, err := a.repo.GetClosedTasksFromTime(today)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	assignedTasksFee := 0
	for _, task := range assignedTasks {
		assignedTasksFee += int(task.PriceForAssign)
	}

	closedTasksCost := 0
	for _, task := range closedTasks {
		closedTasksCost += int(task.PriceForClosing)
	}

	todayEarnings := -1 * (closedTasksCost + assignedTasksFee)

	context.JSON(http.StatusOK, gin.H{"today_earnings": todayEarnings})
}
