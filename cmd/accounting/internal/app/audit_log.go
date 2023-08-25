package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"log"
	"net/http"
	"time"
	"uber-popug/pkg/types"
	"uber-popug/pkg/types/messages"
	v1 "uber-popug/pkg/types/messages/v1"
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

	err := a.repo.CreateAuditLog(auditLog)
	if err != nil {
		return err
	}

	a.produceAnalyticsEvent(auditLog, "closed")

	return nil
}

func (a *App) CreateTaskAssignedAuditLog(task *types.Task) error {
	id, _ := uuid.GenerateUUID()

	auditLog := &types.AuditLog{
		ID:     id,
		UserID: task.AssigneeId,
		Amount: -int(task.PriceForAssign),
		Reason: assignedTask,
		TaskInfo: &types.TaskInfo{
			ID:          task.ID,
			Title:       task.Title,
			JiraID:      task.JiraID,
			Description: task.Description,
		},
		CreatedAt: time.Now(),
	}

	err := a.repo.CreateAuditLog(auditLog)
	if err != nil {
		return err
	}

	a.produceAnalyticsEvent(auditLog, "assigned")

	return nil
}

func (a *App) produceAnalyticsEvent(auditLog *types.AuditLog, status string) {
	id, _ := uuid.GenerateUUID()
	msg := v1.TransactionMessage{
		ID: id,
		Data: v1.TransactionData{
			AuditLogID: auditLog.ID,
			UserID:     auditLog.UserID,
			Amount:     auditLog.Amount,
			Description: fmt.Sprintf("Task %s (JiraID=%s) is %s",
				auditLog.TaskInfo.Title,
				auditLog.TaskInfo.JiraID,
				status),
			CreatedAt: auditLog.CreatedAt,
		},
		CreatedAt: time.Now(),
	}

	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("error producing message")
	}
	a.beProducer.Send(string(res), map[string]string{messages.V1: messages.V1})
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

	logs, err := a.repo.GetUserAuditLogsForPeriod(userID, from, to)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, ToResp(logs))
}

type Response struct {
	AuditLogs []AuditLog `json:"audit_logs"`
}

type AuditLog struct {
	Amount    int       `json:"amount"`
	Reason    string    `json:"reason"`
	TaskInfo  *TaskInfo `json:"task_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskInfo struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	JiraID      string `json:"jira_id,omitempty"`
	Description string `json:"description,omitempty"`
}

func ToResp(logs []*types.AuditLog) Response {
	res := make([]AuditLog, 0, len(logs))

	for _, log := range logs {
		elem := AuditLog{
			Amount:    log.Amount,
			Reason:    log.Reason,
			CreatedAt: log.CreatedAt,
		}

		if log.TaskInfo != nil {
			elem.TaskInfo = &TaskInfo{
				ID:          log.TaskInfo.ID,
				Title:       log.TaskInfo.Title,
				JiraID:      log.TaskInfo.JiraID,
				Description: log.TaskInfo.Description,
			}
		}

		res = append(res, elem)
	}

	return Response{AuditLogs: res}
}
