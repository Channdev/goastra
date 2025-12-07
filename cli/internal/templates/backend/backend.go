package backend

import "fmt"

func GoMod(projectName, db string) string {
	dbImport := "github.com/lib/pq v1.10.9"
	if db == "mysql" {
		dbImport = "github.com/go-sql-driver/mysql v1.7.1"
	}

	return fmt.Sprintf(`module github.com/%s/app

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.16.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	%s
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.16.0
)
`, projectName, dbImport)
}

func MainGo() string {
	return `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if env == "production" {
		godotenv.Load(".env")
	} else {
		godotenv.Load("../../.env." + env)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "version": "1.0.0"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handleLogin)
			auth.POST("/register", handleRegister)
			auth.POST("/refresh", handleRefresh)
			auth.POST("/logout", handleLogout)
		}

		users := v1.Group("/users")
		{
			users.GET("", handleListUsers)
			users.GET("/:id", handleGetUser)
			users.PUT("/:id", handleUpdateUser)
			users.DELETE("/:id", handleDeleteUser)
		}
	}

	if env == "production" {
		staticDir := "./public/browser"
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			filePath := filepath.Join(staticDir, path)
			if _, err := os.Stat(filePath); err == nil {
				c.File(filePath)
				return
			}
			c.File(filepath.Join(staticDir, "index.html"))
		})
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func handleLogin(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login endpoint - implement me"})
}

func handleRegister(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Register endpoint - implement me"})
}

func handleRefresh(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Refresh endpoint - implement me"})
}

func handleLogout(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Logout endpoint - implement me"})
}

func handleListUsers(c *gin.Context) {
	c.JSON(200, gin.H{"data": []interface{}{}, "total": 0})
}

func handleGetUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id")})
}

func handleUpdateUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "updated": true})
}

func handleDeleteUser(c *gin.Context) {
	c.JSON(200, gin.H{"id": c.Param("id"), "deleted": true})
}
`
}
