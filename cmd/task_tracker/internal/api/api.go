package api

import (
	"github.com/gin-gonic/gin"
	"uber-popug/cmd/task_tracker/internal/app"
	"uber-popug/pkg/middlewares"
)

func NewApi(app *app.App) *gin.Engine {
	router := gin.Default()
	group := router.Group("/api").Use(middlewares.Auth())
	{
		group.POST("/task", app.CreateTask)
		group.GET("/tasks/user", app.GetUserTasks)
		group.PATCH("/task", app.CloseTask)
	}

	admin := router.Group("/admin").Use(middlewares.AdminAuth())
	{
		admin.PATCH("/tasks/reassign", app.ReassignTasks)
		admin.DELETE("/task", app.DeleteTask)
	}

	analytics := router.Group("/analytics").Use(middlewares.AdminAuth())
	{
		analytics.GET("/tasks/top", app.TopTask)
		analytics.GET("/today", app.TodayEarnings)
	}

	return router
}
