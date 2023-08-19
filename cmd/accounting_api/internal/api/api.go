package api

import (
	"github.com/gin-gonic/gin"
	"uber-popug/cmd/task_tracker/internal/app"
	"uber-popug/pkg/middlewares"
)

func NewApi(app *app.App) *gin.Engine {
	router := gin.Default()
	admin := router.Group("/admin").Use(middlewares.AdminAuth())
	{
		admin.GET("/analytics/negative", app.NegativePopugs)
		admin.GET("/analytics/tasks/top", app.GetTopTaskOfPeriod) // week/month/day/ in args

		admin.GET("/accounting/balance", app.TotalTodayBalance)
		admin.GET("/accounting/balance/days", app.SeveralDaysStats) // amount of days in args
	}

	popug := router.Group("/accounting").Use(middlewares.Auth())
	{
		popug.GET("/balance", app.PopugBalance)
		popug.GET("/log", app.PopugsTodayLog)
	}

	return router
}
