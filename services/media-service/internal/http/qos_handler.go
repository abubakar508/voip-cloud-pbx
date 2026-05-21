package http

import (
	"net/http"

	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/qos"
	"github.com/gin-gonic/gin"
)

type QoSHandler struct {
	qm *qos.Manager
}

func NewQoSHandler(qm *qos.Manager) *QoSHandler {
	return &QoSHandler{
		qm: qm,
	}
}

func (h *QoSHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/qos", h.handleQoS)
}

func (h *QoSHandler) handleQoS(c *gin.Context) {
	snapshot := h.qm.Snapshot()

	// Convert keys to serializable form
	type StreamKeyDTO struct {
		SSRC uint32 `json:"ssrc"`
		Addr string `json:"addr"`
	}
	type StreamQoS struct {
		Key   StreamKeyDTO    `json:"key"`
		Stats qos.StreamStats `json:"stats"`
	}

	var result []StreamQoS
	for key, stats := range snapshot {
		result = append(result, StreamQoS{
			Key: StreamKeyDTO{
				SSRC: key.SSRC,
				Addr: key.Addr,
			},
			Stats: stats,
		})
	}

	c.JSON(http.StatusOK, result)
}
