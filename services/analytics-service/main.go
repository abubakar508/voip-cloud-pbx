package main

import (
	"fmt"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/analytics-service/internal/events"
	"github.com/abubakar508/voip-cloud-pbx/services/analytics-service/internal/models"
	"github.com/abubakar508/voip-cloud-pbx/services/analytics-service/internal/natsclient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bootstrap := shared.Init("analytics-service")
	shared.PrintBanner("analytics-service")

	logger := bootstrap.Logger

	// DB
	db, err := gorm.Open(postgres.Open(bootstrap.Config.PostgresDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := db.AutoMigrate(&models.CallAnalytics{}); err != nil {
		logger.Fatal("failed to auto-migrate analytics models", zap.Error(err))
	}

	// HTTP
	addr := ":8086"
	server := httpserver.New(httpserver.Options{
		Addr:   addr,
		Logger: logger,
	})
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("analytics http server error", zap.Error(err))
		}
	}()

	// NATS
	nc, err := natsclient.New(bootstrap.Config)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}
	defer nc.Close()
	logger.Info("analytics-service connected to NATS", zap.String("url", nc.URL()))

	// Basic processing: on calls.started store initial row; on calls.ended update duration.
	if _, err := nc.Subscribe("calls.started", func(msg *nats.Msg) {
		var evt events.CallStartedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.started", zap.Error(err))
			return
		}
		rec := models.CallAnalytics{
			CallID:      evt.CallID,
			TenantID:    evt.TenantID,
			FromUser:    evt.FromUser,
			ToUser:      evt.ToUser,
			StartedAt:   evt.StartedAt,
			EndedAt:     evt.StartedAt,
			DurationSec: 0,
		}
		if err := db.Create(&rec).Error; err != nil {
			logger.Error("failed to insert call analytics (start)", zap.Error(err))
		}
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.started", zap.Error(err))
	}

	if _, err := nc.Subscribe("calls.ended", func(msg *nats.Msg) {
		var evt events.CallEndedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.ended", zap.Error(err))
			return
		}
		var rec models.CallAnalytics
		if err := db.Where("call_id = ?", evt.CallID).First(&rec).Error; err != nil {
			logger.Error("failed to find call analytics on end", zap.Error(err))
			return
		}
		rec.EndedAt = evt.EndedAt
		rec.DurationSec = int(evt.EndedAt.Sub(rec.StartedAt).Seconds())
		if err := db.Save(&rec).Error; err != nil {
			logger.Error("failed to update call analytics (end)", zap.Error(err))
		}
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.ended", zap.Error(err))
	}

	fmt.Println("Analytics service started on", addr)

	select {}
}
