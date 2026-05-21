package main

import (
	"fmt"
	"os"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/auth"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/db"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/http"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	bootstrap := shared.Init("auth-service")
	shared.PrintBanner("auth-service")

	// Connect to Postgres
	gdb, err := db.Connect(bootstrap.Config)
	if err != nil {
		bootstrap.Logger.Fatal("failed to connect to postgres", zap.Error(err))
	}

	// Auto-migrate Tenant and User tables
	if err := gdb.AutoMigrate(&models.Tenant{}, &models.User{}); err != nil {
		bootstrap.Logger.Fatal("failed to auto-migrate auth models", zap.Error(err))
	}

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: bootstrap.Logger,
	})

	engine := server.Engine()

	// Basic ping
	engine.GET("/auth/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "auth service ok",
			"port":    port,
		})
	})

	// Register auth routes
	authSvc := auth.NewService(gdb, bootstrap.Config)
	handler := http.NewHandler(authSvc)
	handler.RegisterRoutes(engine)

	fmt.Println("Auth service HTTP server starting on", addr)
	if err := server.Start(); err != nil {
		bootstrap.Logger.Fatal("auth service stopped with error", zap.Error(err))
	}
}
