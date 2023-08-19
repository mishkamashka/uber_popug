package api

import (
	"github.com/gin-gonic/gin"
	"uber-popug/cmd/accounting_api/internal/app"
	"uber-popug/pkg/middlewares"
)

func NewApi(app *app.App) *gin.Engine {
	router := gin.Default()
	admin := router.Group("/admin").Use(middlewares.AdminAuth())
	{
		admin.GET("/analytics/negative", app.NegativePopugs)
	}

	popug := router.Group("/accounting").Use(middlewares.Auth())
	{
		popug.GET("/balance", app.GetPopugBalance)
		popug.GET("/log", app.PopugsTodayLog)
	}

	internal := router.Group("/internal").Use(middlewares.AdminAuth())
	{
		admin.PATCH("/checkout", app.FinalizeDay)
	}

	return router
}
