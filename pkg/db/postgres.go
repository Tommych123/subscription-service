package db

import (
	"fmt"
	"github.com/Tommych123/subscription-service/service/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewPostgres(cfg *config.Config, log *zap.Logger) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB connection", zap.Error(err))
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB", zap.Error(err))
	}
	log.Info("Connected to PostgreSQL database successfully")
	return db
}
