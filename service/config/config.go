package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
}

func LoadConfig(log *zap.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found, loading environment variables directly")
	} else {
		log.Info(".env file loaded successfully")
	}
	cfg := &Config{
		DBHost:     getEnv(log, "DB_HOST", "localhost"),
		DBPort:     getEnv(log, "DB_PORT", "5432"),
		DBUser:     getEnv(log, "DB_USER", "postgres"),
		DBPassword: getEnv(log, "DB_PASSWORD", "postgres"),
		DBName:     getEnv(log, "DB_NAME", "subscriptions"),
		DBSSLMode:  getEnv(log, "DB_SSLMODE", "disable"),
		ServerPort: getEnv(log, "SERVER_PORT", "8080"),
	}
	log.Info("Config loaded",
		zap.String("DBHost", cfg.DBHost),
		zap.String("DBPort", cfg.DBPort),
		zap.String("DBUser", cfg.DBUser),
		zap.String("DBName", cfg.DBName),
		zap.String("DBSSLMode", cfg.DBSSLMode),
		zap.String("ServerPort", cfg.ServerPort),
	)
	return cfg
}

func getEnv(log *zap.Logger, key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Warn("Environment variable not set, using default", zap.String("key", key), zap.String("default", fallback))
	return fallback
}
