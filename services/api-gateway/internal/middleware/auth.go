package middleware

import (
	"net/http"

	"github.com/abubakar508/voip-cloud-pbx/services/api-gateway/internal/authjwt"
	"github.com/gin-gonic/gin"
)

func JWTAuth(v *authjwt.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := v.ParseFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("tenantId", claims.TenantID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
