package db

import (
	"log"
	"time"

	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.AppConfig) (*gorm.DB, error) {
	if cfg.PostgresDSN == "" {
		return nil, ErrMissingDSN
	}

	gormLogger := logger.New(
		log.New(log.Writer(), "gorm: ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
