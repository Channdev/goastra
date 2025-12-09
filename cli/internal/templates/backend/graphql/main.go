/*
 * GoAstra CLI - GraphQL Main Template
 *
 * Generates the main.go entry point for GraphQL servers.
 * Uses gqlgen with Gin for HTTP handling and playground support.
 */
package graphql

// MainGo returns the main.go template for GraphQL servers.
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

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"app/graph"
	"app/graph/generated"
	"app/internal/config"
	"app/internal/database"
	"app/internal/logger"
	"app/internal/middleware"
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

	log.Infow("Starting GraphQL server",
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

	// Initialize GraphQL resolver
	resolver := graph.NewResolver(db)

	// Create GraphQL handler
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	// Setup Gin router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.RequestID())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// GraphQL endpoint
	r.POST("/graphql", gin.WrapH(srv))

	// GraphQL Playground (disabled in production)
	if !cfg.IsProduction() {
		r.GET("/playground", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))
		log.Info("GraphQL Playground available at /playground")
	}

	// Create HTTP server
	httpSrv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Infow("Server listening", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Errorw("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}
`
}
