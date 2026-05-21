package shared

import (
	"fmt"

	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/logger"
	"go.uber.org/zap"
)

type Bootstrap struct {
	Config *config.AppConfig
	Logger *zap.Logger
}

func Init(serviceName string) *Bootstrap {
	cfg := config.Load(serviceName)
	log := logger.New(serviceName, cfg.Env, cfg.LogLevel)

	log.Info("service bootstrap complete",
		zap.String("env", cfg.Env),
		zap.String("version", cfg.Version),
	)

	return &Bootstrap{
		Config: cfg,
		Logger: log,
	}
}

func Version() string {
	return "voip-cloud-pbx-shared-go-0.2.0"
}

func PrintBanner(service string) {
	fmt.Printf("[%s] using shared-go version %s\n", service, Version())
}
