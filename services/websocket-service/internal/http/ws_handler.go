package http

import (
	"net/http"

	"github.com/abubakar508/voip-cloud-pbx/services/websocket-service/internal/authjwt"
	"github.com/abubakar508/voip-cloud-pbx/services/websocket-service/internal/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, restrict origins
		return true
	},
}

type Handler struct {
	hub       *hub.Hub
	validator *authjwt.Validator
}

func NewHandler(h *hub.Hub, v *authjwt.Validator) *Handler {
	return &Handler{
		hub:       h,
		validator: v,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws", h.handleWebSocket)
	r.GET("/ws/ping", h.handlePing)
}

func (h *Handler) handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":        "websocket service ok",
		"connectedCount": h.hub.Count(),
	})
}

func (h *Handler) handleWebSocket(c *gin.Context) {
	token := c.Query("token")
	authHeader := c.GetHeader("Authorization")

	claims, err := h.validator.ParseFromQueryOrHeader(token, authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &hub.Client{
		Conn:     conn,
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
	}

	h.hub.Add(client)

	go func() {
		defer func() {
			h.hub.Remove(conn)
			_ = conn.Close()
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// For now broadcast received message to all
			h.hub.Broadcast(messageType, message)
		}
	}()
}
