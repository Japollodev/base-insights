package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypto-insights/internal/config"
	"crypto-insights/internal/database"
	"crypto-insights/internal/logger"
	"crypto-insights/internal/server"
	"crypto-insights/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redis, err := database.NewRedis(cfg.Redis)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize services
	blockchainService := services.NewBlockchainService(cfg.Blockchain, logger)
	dataCollectionService := services.NewDataCollectionService(cfg.DataCollection, blockchainService, db, redis, logger)

	// Initialize HTTP server
	server := server.New(cfg.Server.Port, blockchainService, dataCollectionService, logger)

	// Start data collection service
	if err := dataCollectionService.Start(context.Background()); err != nil {
		logger.Fatalf("Failed to start data collection service: %v", err)
	}

	// Start HTTP server
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	if err := dataCollectionService.Stop(ctx); err != nil {
		logger.Errorf("Failed to stop data collection service: %v", err)
	}

	logger.Info("Server exited")
}
