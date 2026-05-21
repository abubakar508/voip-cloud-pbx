package http

import (
	"net/http"
	"time"

	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/calls"
	"github.com/gin-gonic/gin"
)

type CallsListHandler struct {
	callMgr *calls.Manager
}

func NewCallsListHandler(callMgr *calls.Manager) *CallsListHandler {
	return &CallsListHandler{
		callMgr: callMgr,
	}
}

func (h *CallsListHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/calls", h.handleListCalls)
}

func (h *CallsListHandler) handleListCalls(c *gin.Context) {
	// Simple snapshot: collect sessions from manager
	type CallDTO struct {
		CallID    string `json:"callId"`
		TenantID  string `json:"tenantId"`
		FromUser  string `json:"fromUser"`
		ToUser    string `json:"toUser"`
		Direction string `json:"direction"`
		StartedAt string `json:"startedAt"`
		EndedAt   string `json:"endedAt,omitempty"`
	}

	var result []CallDTO

	// There is no direct snapshot method, so we can add one or do this:
	// We'll add a Snapshot method to Manager.

	callsSnapshot := h.callMgr.Snapshot()

	for _, s := range callsSnapshot {
		dto := CallDTO{
			CallID:    s.CallID,
			TenantID:  s.TenantID,
			FromUser:  s.FromUser,
			ToUser:    s.ToUser,
			Direction: string(s.Direction),
			StartedAt: s.StartedAt.UTC().Format(time.RFC3339),
		}
		if s.EndedAt != nil {
			dto.EndedAt = s.EndedAt.UTC().Format(time.RFC3339)
		}
		result = append(result, dto)
	}

	c.JSON(http.StatusOK, result)
}
