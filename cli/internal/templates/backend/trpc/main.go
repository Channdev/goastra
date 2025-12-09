/*
 * GoAstra CLI - tRPC Main Template
 *
 * Generates the main.go entry point for tRPC servers.
 * Uses Connect-Go for type-safe RPC with HTTP support.
 */
package trpc

// MainGo returns the main.go template for tRPC servers.
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

	"connectrpc.com/connect"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"app/internal/config"
	"app/internal/database"
	"app/internal/logger"
	"app/internal/rpc"
	"app/internal/rpc/gen/proto/v1/protov1connect"
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

	log.Infow("Starting tRPC server",
		"env", cfg.Env,
		"port", cfg.Port,
	)

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

	// Create HTTP mux
	mux := http.NewServeMux()

	// Create interceptor for logging and auth
	interceptor := connect.WithInterceptors(
		rpc.NewLoggingInterceptor(log),
	)

	// Register services
	userService := rpc.NewUserService(db, cfg)
	authService := rpc.NewAuthService(db, cfg)
	healthService := rpc.NewHealthService(db)

	// Mount service handlers
	path, handler := protov1connect.NewUserServiceHandler(userService, interceptor)
	mux.Handle(path, handler)

	path, handler = protov1connect.NewAuthServiceHandler(authService, interceptor)
	mux.Handle(path, handler)

	path, handler = protov1connect.NewHealthServiceHandler(healthService, interceptor)
	mux.Handle(path, handler)

	// CORS handler
	corsHandler := corsMiddleware(cfg.CORSOrigins, mux)

	// Create HTTP server with HTTP/2 support
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      h2c.NewHandler(corsHandler, &http2.Server{}),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Infow("Server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorw("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}

func corsMiddleware(allowedOrigins string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
`
}
