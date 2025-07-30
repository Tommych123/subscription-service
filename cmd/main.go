// @title Subscription Service API
// @version 1.0
// @description REST API для управления подписками пользователей
// @host localhost:8080
// @BasePath /
package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"github.com/Tommych123/subscription-service/service/config"
	"github.com/Tommych123/subscription-service/pkg/db"
	"github.com/Tommych123/subscription-service/api"
	"github.com/Tommych123/subscription-service/internal/logger"
	"github.com/Tommych123/subscription-service/repository"
	"github.com/Tommych123/subscription-service/service"
	_ "github.com/Tommych123/subscription-service/internal/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func main() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found, reading from environment")
	}
	logg, err := logger.NewLogger()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logg.Sync()
	logg.Info("Starting application")
	cfg := config.LoadConfig(logg)
	sqlxDB := db.NewPostgres(cfg, logg)
	db.RunMigrations(sqlxDB, cfg, logg)
	repo := repository.NewSubscriptionRepository(sqlxDB)
	svc := service.NewSubscriptionService(repo, logg)
	h := api.NewSubscriptionHandler(svc, logg)
	r := gin.Default()
	h.RegisterRoutes(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = cfg.ServerPort
	}
	logg.Info("Server running", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logg.Fatal("Server failed", zap.Error(err))
	}
}
