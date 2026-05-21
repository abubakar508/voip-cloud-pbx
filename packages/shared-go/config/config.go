package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Name     string
	Version  string
	Env      string
	LogLevel string

	HTTPPort string

	PostgresDSN string
	RedisAddr   string
	NATSURL     string

	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration

	RedisPassword string
	RedisHost     string
	RedisPort     string
	NatsURL       string
}

func Load(serviceName string) *AppConfig {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("LOG_LEVEL", "info")

	cfg := &AppConfig{
		Name:     serviceName,
		Version:  "0.1.0",
		Env:      v.GetString("APP_ENV"),
		LogLevel: v.GetString("LOG_LEVEL"),

		HTTPPort: v.GetString("PORT"),

		PostgresDSN: buildPostgresDSN(v),
		RedisAddr:   buildRedisAddr(v),
		NATSURL:     v.GetString("NATS_URL"),

		JWTAccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
	}

	accessTTLStr := v.GetString("JWT_ACCESS_TTL")
	refreshTTLStr := v.GetString("JWT_REFRESH_TTL")

	if accessTTLStr == "" {
		accessTTLStr = "15m"
	}
	if refreshTTLStr == "" {
		refreshTTLStr = "720h"
	}

	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		log.Printf("invalid JWT_ACCESS_TTL, using default 15m: %v", err)
		accessTTL = 15 * time.Minute
	}
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		log.Printf("invalid JWT_REFRESH_TTL, using default 720h: %v", err)
		refreshTTL = 30 * 24 * time.Hour
	}

	cfg.JWTAccessTTL = accessTTL
	cfg.JWTRefreshTTL = refreshTTL

	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}

	return cfg
}

func buildPostgresDSN(v *viper.Viper) string {
	host := v.GetString("POSTGRES_HOST")
	port := v.GetString("POSTGRES_PORT")
	user := v.GetString("POSTGRES_USER")
	pass := v.GetString("POSTGRES_PASSWORD")
	db := v.GetString("POSTGRES_DB")

	if host == "" || user == "" || db == "" {
		return ""
	}

	if port == "" {
		port = "5432"
	}

	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=disable"
}

func buildRedisAddr(v *viper.Viper) string {
	host := v.GetString("REDIS_HOST")
	port := v.GetString("REDIS_PORT")
	if host == "" {
		return ""
	}
	if port == "" {
		port = "6379"
	}
	return host + ":" + port
}
