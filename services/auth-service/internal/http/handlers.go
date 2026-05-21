package http

import (
	"net/http"

	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authSvc *auth.Service
}

func NewHandler(authSvc *auth.Service) *Handler {
	return &Handler{authSvc: authSvc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/auth/register", h.handleRegister)
	r.POST("/auth/login", h.handleLogin)
}

func (h *Handler) handleRegister(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	user, err := h.authSvc.RegisterTenantAdmin(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId":   user.ID,
		"tenantId": user.TenantID,
		"email":    user.Email,
	})
}

func (h *Handler) handleLogin(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	user, access, refresh, err := h.authSvc.Login(req)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":       user.ID,
		"tenantId":     user.TenantID,
		"accessToken":  access,
		"refreshToken": refresh,
	})
}
