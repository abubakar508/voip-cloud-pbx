package main

import (
	"fmt"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/ai-service/internal/events"
	"github.com/abubakar508/voip-cloud-pbx/services/ai-service/internal/models"
	"github.com/abubakar508/voip-cloud-pbx/services/ai-service/internal/natsclient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bootstrap := shared.Init("ai-service")
	shared.PrintBanner("ai-service")

	logger := bootstrap.Logger

	// DB
	db, err := gorm.Open(postgres.Open(bootstrap.Config.PostgresDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := db.AutoMigrate(&models.AISummary{}); err != nil {
		logger.Fatal("failed to auto-migrate ai models", zap.Error(err))
	}

	// HTTP
	addr := ":8085"
	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: logger,
	})
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("ai http server error", zap.Error(err))
		}
	}()

	// NATS
	nc, err := natsclient.New(bootstrap.Config)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}
	defer nc.Close()
	logger.Info("ai-service connected to NATS", zap.String("url", nc.URL()))

	// Subscribe to calls.ended (future: trigger AI processing)
	if _, err := nc.Subscribe("calls.ended", func(msg *nats.Msg) {
		var evt events.CallEndedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.ended", zap.Error(err))
			return
		}
		logger.Info("ai-service saw calls.ended",
			zap.String("callId", evt.CallID),
			zap.String("tenantId", evt.TenantID),
			zap.String("reason", evt.Reason),
		)
		// In future: load transcript/recording, run AI model, store AISummary in DB.
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.ended", zap.Error(err))
	}

	fmt.Println("AI service started on", addr)

	select {}
}
