package http

import (
	"net/http"
	"strings"

	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserKey = "authUser"
)

func AuthMiddleware(authSvc *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := authSvc.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		ctxUser := &auth.ContextUser{
			ID:       claims.UserID,
			TenantID: claims.TenantID,
			Role:     claims.Role,
			Email:    "", // can be filled by querying user if needed
		}

		c.Set(ContextUserKey, ctxUser)
		c.Next()
	}
}

func GetContextUser(c *gin.Context) (*auth.ContextUser, bool) {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*auth.ContextUser)
	return u, ok
}
