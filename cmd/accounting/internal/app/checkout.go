package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"

	"uber-popug/pkg/accounting"
	"uber-popug/pkg/types"
)

func (a *App) Checkout(context *gin.Context) {
	req := accounting.CheckoutRequest{}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	err := a.repo.UpdatePopugBalanceByValue(req.UserID, -req.DayTotal)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	id, _ := uuid.GenerateUUID()
	err = a.repo.CreateAuditLog(&types.AuditLog{
		ID:     id,
		UserID: req.UserID,
		Amount: -req.DayTotal,
		Reason: "daily checkout",
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
	}
}
