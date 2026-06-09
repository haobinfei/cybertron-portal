package main

import (
	"fmt"
	"log"

	"cybertron-portal/internal/config"
	"cybertron-portal/internal/handler"
	"cybertron-portal/internal/model"
	"cybertron-portal/internal/repository"
	"cybertron-portal/internal/router"
	"cybertron-portal/internal/service"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Fatal("Failed to auto migrate", zap.Error(err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, rdb, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	userService := service.NewUserService(userRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(authService, authHandler, userHandler)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
