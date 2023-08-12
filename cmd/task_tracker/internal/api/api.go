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
		group.GET("/task", app.GetTask)
		group.DELETE("/task", app.DeleteTask)
		group.POST("/task/status", app.UpdateTask)

		group.GET("/user/tasks", app.GetUserTasks)
	}

	admin := router.Group("/admin").Use(middlewares.AdminAuth())
	{
		admin.POST("/tasks/reassign", app.ReassignTasks)
	}

	return router
}