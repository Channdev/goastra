/*
 * GoAstra Backend - Main Entry Point
 *
 * Initializes the application, loads configuration,
 * establishes database connection, and starts the HTTP server
 * with graceful shutdown support.
 */
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/channdev/goastra/app/internal/config"
	"github.com/channdev/goastra/app/internal/database"
	"github.com/channdev/goastra/app/internal/logger"
	"github.com/channdev/goastra/app/internal/router"
)

func main() {
	/* Determine environment and load corresponding .env file */
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env
	if err := godotenv.Load("../../" + envFile); err != nil {
		log.Printf("No %s file found, using environment variables", envFile)
	}

	/* Initialize configuration from environment */
	cfg := config.Load()

	/* Initialize structured logger */
	appLogger := logger.New(cfg.LogLevel)
	defer appLogger.Sync()

	/* Establish database connection */
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	/* Configure HTTP router with middleware and routes */
	r := router.New(appLogger, db, cfg)

	/* Configure HTTP server with timeouts */
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	/* Start server in goroutine for graceful shutdown */
	go func() {
		appLogger.Info("Server starting", "port", cfg.Port, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed to start", "error", err)
		}
	}()

	/* Wait for interrupt signal */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	/* Graceful shutdown with timeout */
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited gracefully")
}
