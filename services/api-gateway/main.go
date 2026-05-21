package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/api-gateway/internal/authjwt"
	"github.com/abubakar508/voip-cloud-pbx/services/api-gateway/internal/middleware"
	"github.com/abubakar508/voip-cloud-pbx/services/api-gateway/internal/proxy"
)

func main() {
	bootstrap := shared.Init("api-gateway")
	shared.PrintBanner("api-gateway")

	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: bootstrap.Logger,
	})

	engine := server.Engine()

	engine.GET("/gateway/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "api gateway ok"})
	})

	validator := authjwt.NewValidator(bootstrap.Config)

	// Public auth routes proxy to auth-service
	authBase := "http://voip-auth-service:8081"
	engine.Any("/auth/*path", proxy.ProxyHandler(authBase))

	// Protected example route (would proxy to other services later)
	protected := engine.Group("/api")
	protected.Use(middleware.JWTAuth(validator))
	protected.GET("/profile", func(c *gin.Context) {
		userId, _ := c.Get("userId")
		tenantId, _ := c.Get("tenantId")
		role, _ := c.Get("role")
		c.JSON(200, gin.H{
			"userId":   userId,
			"tenantId": tenantId,
			"role":     role,
		})
	})

	fmt.Println("API Gateway HTTP server starting on", addr)
	if err := server.Start(); err != nil {
		bootstrap.Logger.Fatal("api gateway stopped with error", zap.Error(err))
	}
}