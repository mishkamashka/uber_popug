package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/util"
)

const (
	closedTask   = "task closed"
	assignedTask = "task assigned"
)

func (a *App) CreateTaskClosedAuditLog(task *types.Task) error {
	id, _ := uuid.GenerateUUID()

	auditLog := &types.AuditLog{
		ID:     id,
		UserID: task.AssigneeId,
		Amount: int(task.PriceForClosing),
		Reason: closedTask,
		TaskInfo: &types.TaskInfo{
			ID:          task.ID,
			Title:       task.Title,
			JiraID:      task.JiraID,
			Description: task.Description,
		},
	}

	return a.repo.CreateAuditLog(auditLog)
}

func (a *App) CreateTaskAssignedAuditLog(task *types.Task) error {
	id, _ := uuid.GenerateUUID()

	auditLog := &types.AuditLog{
		ID:     id,
		UserID: task.AssigneeId,
		Amount: -int(task.PriceForClosing),
		Reason: assignedTask,
		TaskInfo: &types.TaskInfo{
			ID:          task.ID,
			Title:       task.Title,
			JiraID:      task.JiraID,
			Description: task.Description,
		},
	}

	return a.repo.CreateAuditLog(auditLog)
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
