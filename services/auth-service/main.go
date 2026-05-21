package main

import (
	"fmt"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/gin-gonic/gin"
)

func main() {
	bootstrap := shared.Init("auth-service")
	shared.PrintBanner("auth-service")

	addr := ":" + bootstrap.Config.HTTPPort
	if bootstrap.Config.HTTPPort == "" {
		addr = ":8081"
	}

	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: bootstrap.Logger,
	})

	// add an auth-specific health/info endpoint for now
	engine := server.Engine()
	engine.GET("/auth/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "auth service ok"})
	})

	fmt.Println("Auth service HTTP server starting on", addr)
	if err := server.Start(); err != nil {
		bootstrap.Logger.Fatal("auth service stopped with error", zapError(err))
	}
}

func zapError(err error) interface{} {
	if err == nil {
		return nil
	}
	return err
}
