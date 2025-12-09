/*
 * GoAstra CLI - REST Main Template
 *
 * Generates the main.go entry point for REST API servers.
 * Uses Gin framework with graceful shutdown support.
 */
package rest

// MainGo returns the main.go template for REST API servers.
func MainGo() string {
	return `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"app/internal/config"
	"app/internal/database"
	"app/internal/logger"
	"app/internal/router"
)

func main() {
	// Load environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Load .env file
	if env == "production" {
		godotenv.Load(".env")
	} else {
		godotenv.Load("../../.env." + env)
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.Env)
	defer log.Sync()

	log.Infow("Starting server",
		"env", cfg.Env,
		"port", cfg.Port,
	)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalw("Failed to connect to database", "error", err)
	}
	if db != nil {
		defer db.Close()
		log.Info("Database connected")
	} else {
		log.Warn("No database configured")
	}

	// Initialize router
	r := router.New(log, db, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Infow("Server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorw("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}
`
}
