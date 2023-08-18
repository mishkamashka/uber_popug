package api

import (
	"github.com/gin-gonic/gin"
	"uber-popug/cmd/task_tracker/internal/app"
	"uber-popug/pkg/middlewares"
)

func NewApi(app *app.App) *gin.Engine {
	router := gin.Default()
	group := router.Group("/admin").Use(middlewares.AdminAuth())
	{
		group.GET("/analytics/today", app.TodayBalance)
		group.GET("/analytics/negative", app.NegativePopugs)
		group.GET("/analytics/tasks/top", app.GetTopTaskOfPeriod) // week/month/day/ in args

		group.GET("/accounting/balance", app.TodayBalance)
		group.GET("/accounting/balance/days", app.SeveralDaysStats) // amount of days in args
	}

	admin := router.Group("/accounting").Use(middlewares.Auth())
	{
		group.GET("/balance", app.PopugBalance)
		group.GET("/log", app.PopugsTodayLog)
	}

	return router
}
