package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *App) GetPopugBalance(context *gin.Context) {
	userID, ok := context.Value("userID").(string)
	if userID == "" || !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no userID or userID is not string"})
		context.Abort()
		return
	}

	balance, err := a.repo.GetPopugBalance(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, balance)
}

func (a *App) UpdatePopugBalance(userID string, amount int) error {
	err := a.repo.UpdatePopugBalanceByValue(userID, amount)
	if err != nil {
		return fmt.Errorf("update popug balance: %s", err)
	}

	return nil
}
