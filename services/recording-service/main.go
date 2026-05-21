package main

import (
	"fmt"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/recording-service/internal/events"
	"github.com/abubakar508/voip-cloud-pbx/services/recording-service/internal/models"
	"github.com/abubakar508/voip-cloud-pbx/services/recording-service/internal/natsclient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bootstrap := shared.Init("recording-service")
	shared.PrintBanner("recording-service")

	logger := bootstrap.Logger

	// DB
	db, err := gorm.Open(postgres.Open(bootstrap.Config.PostgresDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := db.AutoMigrate(&models.Recording{}, &models.CallRecord{}); err != nil {
		logger.Fatal("failed to auto-migrate recording models", zap.Error(err))
	}

	// HTTP server (health only for now)
	addr := ":8083"
	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: logger,
	})
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("recording http server error", zap.Error(err))
		}
	}()

	// NATS
	nc, err := natsclient.New(bootstrap.Config)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}
	defer nc.Close()
	logger.Info("recording-service connected to NATS", zap.String("url", nc.URL()))

	// Subscribe to call events (no heavy logic yet, just log)
	if _, err := nc.Subscribe("calls.started", func(msg *nats.Msg) {
		var evt events.CallStartedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.started", zap.Error(err))
			return
		}
		logger.Info("recording-service saw calls.started",
			zap.String("callId", evt.CallID),
			zap.String("tenantId", evt.TenantID),
		)
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.started", zap.Error(err))
	}

	if _, err := nc.Subscribe("calls.ended", func(msg *nats.Msg) {
		var evt events.CallEndedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.ended", zap.Error(err))
			return
		}
		logger.Info("recording-service saw calls.ended",
			zap.String("callId", evt.CallID),
			zap.String("tenantId", evt.TenantID),
			zap.String("reason", evt.Reason),
		)
		// In future: finalize recording, update DB, etc.
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.ended", zap.Error(err))
	}

	fmt.Println("Recording service started on", addr)

	select {} // block forever
}
