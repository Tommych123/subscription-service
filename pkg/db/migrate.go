package db

import (
	"github.com/Tommych123/subscription-service/service/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunMigrations(db *sqlx.DB, cfg *config.Config, log *zap.Logger) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create DB migration driver", zap.Error(err))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		log.Fatal("Failed to initialize migrate instance", zap.Error(err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed", zap.Error(err))
	}
	log.Info("Migrations applied successfully")
}
