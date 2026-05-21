package http

import (
	"net"
	"net/http"

	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/calls"
	"github.com/gin-gonic/gin"
)

type CallsHandler struct {
	callMgr *calls.Manager
}

func NewCallsHandler(callMgr *calls.Manager) *CallsHandler {
	return &CallsHandler{
		callMgr: callMgr,
	}
}

type SetEndpointsRequest struct {
	AAddr string `json:"aAddr" binding:"required"`
	BAddr string `json:"bAddr" binding:"required"`
}

func (h *CallsHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/calls/:callId/endpoints", h.handleSetEndpoints)
}

func (h *CallsHandler) handleSetEndpoints(c *gin.Context) {
	callID := c.Param("callId")

	var req SetEndpointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	aUDP, err := net.ResolveUDPAddr("udp", req.AAddr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid AAddr"})
		return
	}
	bUDP, err := net.ResolveUDPAddr("udp", req.BAddr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid BAddr"})
		return
	}

	h.callMgr.SetEndpoints(callID, aUDP, bUDP)

	c.JSON(http.StatusOK, gin.H{
		"callId": callID,
		"aAddr":  aUDP.String(),
		"bAddr":  bUDP.String(),
	})
}
