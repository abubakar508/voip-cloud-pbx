package main

import (
	"fmt"
	"os"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/websocket-service/internal/authjwt"
	"github.com/abubakar508/voip-cloud-pbx/services/websocket-service/internal/http"
	"github.com/abubakar508/voip-cloud-pbx/services/websocket-service/internal/hub"
)

func main() {
	bootstrap := shared.Init("websocket-service")
	shared.PrintBanner("websocket-service")

	port := os.Getenv("WEBSOCKET_SERVICE_PORT")
	if port == "" {
		port = "8084"
	}
	addr := ":" + port

	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: bootstrap.Logger,
	})

	engine := server.Engine()

	h := hub.New()
	validator := authjwt.NewValidator(bootstrap.Config)

	handler := http.NewHandler(h, validator)
	handler.RegisterRoutes(engine)

	fmt.Println("Websocket service HTTP server starting on", addr)
	if err := server.Start(); err != nil {
		bootstrap.Logger.Fatal("websocket service stopped with error", zap.Error(err))
	}
}
