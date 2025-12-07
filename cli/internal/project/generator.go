/*
 * GoAstra CLI - Project Generator
 *
 * Handles scaffolding of new GoAstra projects.
 * Creates directory structure, configuration files, and boilerplate code.
 */
package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goastra/cli/internal/env"
)

/*
 * Config holds the project generation configuration.
 * Populated from CLI flags and defaults.
 */
type Config struct {
	Name        string
	Path        string
	SkipAngular bool
	SkipBackend bool
	UseGraphQL  bool
}

/*
 * Generator handles project scaffolding operations.
 * Uses embedded templates for file generation.
 */
type Generator struct {
	config *Config
}

/*
 * NewGenerator creates a new project generator instance.
 */
func NewGenerator(cfg *Config) *Generator {
	return &Generator{config: cfg}
}

/*
 * CreateStructure creates the complete directory hierarchy.
 * Follows the GoAstra monorepo convention.
 */
func (g *Generator) CreateStructure() error {
	dirs := []string{
		"app/cmd/server",
		"app/internal/config",
		"app/internal/middleware",
		"app/internal/handlers",
		"app/internal/models",
		"app/internal/repository",
		"app/internal/services",
		"app/internal/router",
		"app/internal/auth",
		"app/internal/logger",
		"app/internal/database",
		"app/internal/validator",
		"app/migrations",
		"web/src/app/core/services",
		"web/src/app/core/guards",
		"web/src/app/core/interceptors",
		"web/src/app/core/models",
		"web/src/app/shared/components",
		"web/src/app/shared/directives",
		"web/src/app/shared/pipes",
		"web/src/app/features",
		"web/src/environments",
		"web/src/assets",
		"schema/types",
		"schema/templates",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(g.config.Path, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

/*
 * GenerateEnvFiles creates environment configuration files.
 * Creates .env.development, .env.production, and .env.test.
 */
func (g *Generator) GenerateEnvFiles() error {
	envConfigs := map[string]map[string]string{
		".env.development": {
			"APP_ENV":    "development",
			"API_URL":    "http://localhost:8080",
			"DB_URL":     "postgres://user:password@localhost:5432/goastra_dev",
			"JWT_SECRET": "dev-secret-change-in-production",
			"PORT":       "8080",
			"LOG_LEVEL":  "debug",
		},
		".env.production": {
			"APP_ENV":    "production",
			"API_URL":    "https://api.example.com",
			"DB_URL":     "postgres://user:password@localhost:5432/goastra_prod",
			"JWT_SECRET": "",
			"PORT":       "8080",
			"LOG_LEVEL":  "info",
		},
		".env.test": {
			"APP_ENV":    "test",
			"API_URL":    "http://localhost:8081",
			"DB_URL":     "postgres://user:password@localhost:5432/goastra_test",
			"JWT_SECRET": "test-secret",
			"PORT":       "8081",
			"LOG_LEVEL":  "error",
		},
	}

	for filename, vars := range envConfigs {
		path := filepath.Join(g.config.Path, filename)
		if err := env.WriteEnvFile(path, vars); err != nil {
			return fmt.Errorf("failed to create %s: %w", filename, err)
		}
	}

	return nil
}

/*
 * GenerateConfig creates the goastra.json configuration file.
 * Defines project settings and code generation rules.
 */
func (g *Generator) GenerateConfig() error {
	apiType := "rest"
	if g.config.UseGraphQL {
		apiType = "graphql"
	}

	config := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "api": {
    "type": "%s",
    "prefix": "/api/v1"
  },
  "backend": {
    "port": 8080,
    "module": "github.com/%s/app"
  },
  "frontend": {
    "port": 4200,
    "proxy": "/api"
  },
  "codegen": {
    "schemaPath": "schema/types",
    "outputPath": "web/src/app/core/models",
    "serviceOutput": "web/src/app/core/services"
  },
  "database": {
    "driver": "postgres",
    "migrationsPath": "app/migrations"
  }
}`, g.config.Name, apiType, g.config.Name)

	configPath := filepath.Join(g.config.Path, "goastra.json")
	return os.WriteFile(configPath, []byte(config), 0644)
}

/*
 * GenerateBackend creates the Go backend boilerplate.
 * Includes main entry point, configuration, and router setup.
 */
func (g *Generator) GenerateBackend() error {
	if err := g.generateGoMod(); err != nil {
		return err
	}

	if err := g.generateMainGo(); err != nil {
		return err
	}

	if err := g.generateConfigLoader(); err != nil {
		return err
	}

	if err := g.generateRouter(); err != nil {
		return err
	}

	if err := g.generateMiddleware(); err != nil {
		return err
	}

	if err := g.generateLogger(); err != nil {
		return err
	}

	if err := g.generateDatabase(); err != nil {
		return err
	}

	return nil
}

/*
 * GenerateFrontend creates the Angular frontend boilerplate.
 * Includes app module, core module, and environment configuration.
 */
func (g *Generator) GenerateFrontend() error {
	if err := g.generatePackageJSON(); err != nil {
		return err
	}

	if err := g.generateAngularJSON(); err != nil {
		return err
	}

	if err := g.generateTsConfig(); err != nil {
		return err
	}

	if err := g.generateAppModule(); err != nil {
		return err
	}

	if err := g.generateCoreModule(); err != nil {
		return err
	}

	if err := g.generateEnvironments(); err != nil {
		return err
	}

	return nil
}

/*
 * GenerateSchema creates the shared type definitions structure.
 * Sets up Go structs that will be converted to TypeScript.
 */
func (g *Generator) GenerateSchema() error {
	schemaContent := `/*
 * GoAstra Schema Types
 *
 * Define your shared types here. These will be converted to TypeScript
 * using the 'goastra typesync' command.
 */
package types

import "time"

/*
 * BaseModel provides common fields for all database models.
 * Automatically handles ID, timestamps, and soft deletes.
 */
type BaseModel struct {
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

/*
 * User represents an authenticated user in the system.
 */
type User struct {
	BaseModel
	Email    string ` + "`json:\"email\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	Role     string ` + "`json:\"role\"`" + `
	Active   bool   ` + "`json:\"active\"`" + `
}

/*
 * PaginatedResponse wraps list responses with pagination metadata.
 */
type PaginatedResponse struct {
	Data       interface{} ` + "`json:\"data\"`" + `
	Total      int         ` + "`json:\"total\"`" + `
	Page       int         ` + "`json:\"page\"`" + `
	PageSize   int         ` + "`json:\"page_size\"`" + `
	TotalPages int         ` + "`json:\"total_pages\"`" + `
}

/*
 * APIError represents a standardized error response.
 */
type APIError struct {
	Code    string ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
	Details interface{} ` + "`json:\"details,omitempty\"`" + `
}
`

	schemaPath := filepath.Join(g.config.Path, "schema/types/types.go")
	return os.WriteFile(schemaPath, []byte(schemaContent), 0644)
}

func (g *Generator) generateGoMod() error {
	content := fmt.Sprintf(`module github.com/%s/app

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.5.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.26.0
)
`, g.config.Name)

	return os.WriteFile(filepath.Join(g.config.Path, "app/go.mod"), []byte(content), 0644)
}

func (g *Generator) generateMainGo() error {
	content := `/*
 * GoAstra Backend - Main Entry Point
 *
 * Initializes the application, loads configuration,
 * and starts the HTTP server.
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
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env
	if err := godotenv.Load("../" + envFile); err != nil {
		log.Printf("No %s file found, using environment variables", envFile)
	}

	/* Initialize components */
	cfg := LoadConfig()
	logger := NewLogger(cfg.LogLevel)
	db := InitDatabase(cfg.DatabaseURL)
	defer db.Close()

	router := SetupRouter(logger, db)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	/* Graceful shutdown */
	go func() {
		logger.Info("Server starting on port " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exited")
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/main.go"), []byte(content), 0644)
}

func (g *Generator) generateConfigLoader() error {
	content := `/*
 * GoAstra Backend - Configuration Loader
 *
 * Reads configuration from environment variables.
 * Provides type-safe access to all configuration values.
 */
package main

import "os"

/*
 * Config holds all application configuration values.
 */
type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	JWTSecret   string
	LogLevel    string
	APIURL      string
}

/*
 * LoadConfig reads configuration from environment.
 * Falls back to sensible defaults for development.
 */
func LoadConfig() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DB_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/config.go"), []byte(content), 0644)
}

func (g *Generator) generateRouter() error {
	content := `/*
 * GoAstra Backend - Router Setup
 *
 * Configures HTTP routing with Gin framework.
 * Registers middleware and API routes.
 */
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

/*
 * SetupRouter initializes the Gin router with all routes and middleware.
 */
func SetupRouter(logger *Logger, db *sqlx.DB) *gin.Engine {
	router := gin.New()

	/* Global middleware */
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(RequestLogger(logger))

	/* Health check endpoint */
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	/* API v1 routes */
	v1 := router.Group("/api/v1")
	{
		/* Auth routes */
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handleLogin(db))
			auth.POST("/register", handleRegister(db))
			auth.POST("/refresh", handleRefresh())
		}

		/* Protected routes */
		protected := v1.Group("")
		protected.Use(AuthMiddleware())
		{
			/* User routes */
			users := protected.Group("/users")
			{
				users.GET("", handleListUsers(db))
				users.GET("/:id", handleGetUser(db))
				users.PUT("/:id", handleUpdateUser(db))
				users.DELETE("/:id", handleDeleteUser(db))
			}
		}
	}

	return router
}

/* Placeholder handlers - implement in handlers package */
func handleLogin(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleRegister(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleListUsers(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleGetUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleUpdateUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}

func handleDeleteUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{"error": "not implemented"})
	}
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/router.go"), []byte(content), 0644)
}

func (g *Generator) generateMiddleware() error {
	content := `/*
 * GoAstra Backend - Middleware
 *
 * HTTP middleware for authentication, logging, and CORS.
 */
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/*
 * CORSMiddleware handles Cross-Origin Resource Sharing.
 * Configures allowed origins, methods, and headers.
 */
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

/*
 * AuthMiddleware validates JWT tokens and sets user context.
 */
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		/* TODO: Validate JWT and extract claims */
		claims, err := ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

/*
 * RequestLogger logs HTTP requests with timing information.
 */
func RequestLogger(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("request",
			"method", c.Request.Method,
			"path", path,
			"status", statusCode,
			"duration", duration.String(),
			"ip", c.ClientIP(),
		)
	}
}

/*
 * RoleMiddleware restricts access to specific user roles.
 */
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no role found"})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/middleware.go"), []byte(content), 0644)
}

func (g *Generator) generateLogger() error {
	content := `/*
 * GoAstra Backend - Logger
 *
 * Structured logging using zap for high-performance logging.
 */
package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
 * Logger wraps zap.SugaredLogger for convenience.
 */
type Logger struct {
	*zap.SugaredLogger
}

/*
 * NewLogger creates a new logger instance with the specified level.
 */
func NewLogger(level string) *Logger {
	config := zap.NewProductionConfig()

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	return &Logger{logger.Sugar()}
}

/*
 * WithFields returns a logger with additional context fields.
 */
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	return &Logger{l.SugaredLogger.With(fields...)}
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/logger.go"), []byte(content), 0644)
}

func (g *Generator) generateDatabase() error {
	content := `/*
 * GoAstra Backend - Database Connection
 *
 * Handles PostgreSQL connection and pooling using sqlx.
 */
package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
 * InitDatabase establishes connection to PostgreSQL.
 * Configures connection pool settings for production use.
 */
func InitDatabase(databaseURL string) *sqlx.DB {
	if databaseURL == "" {
		log.Println("Warning: No database URL provided, running without database")
		return nil
	}

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	/* Configure connection pool */
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	/* Verify connection */
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")

	return db
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "app/cmd/server/database.go"), []byte(content), 0644)
}

func (g *Generator) generatePackageJSON() error {
	content := fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "build:prod": "ng build --configuration production",
    "test": "ng test",
    "lint": "ng lint"
  },
  "dependencies": {
    "@angular/animations": "^17.0.0",
    "@angular/common": "^17.0.0",
    "@angular/compiler": "^17.0.0",
    "@angular/core": "^17.0.0",
    "@angular/forms": "^17.0.0",
    "@angular/platform-browser": "^17.0.0",
    "@angular/platform-browser-dynamic": "^17.0.0",
    "@angular/router": "^17.0.0",
    "rxjs": "~7.8.0",
    "tslib": "^2.6.0",
    "zone.js": "~0.14.0"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^17.0.0",
    "@angular/cli": "^17.0.0",
    "@angular/compiler-cli": "^17.0.0",
    "@types/jasmine": "~5.1.0",
    "@types/node": "^20.0.0",
    "jasmine-core": "~5.1.0",
    "karma": "~6.4.0",
    "karma-chrome-launcher": "~3.2.0",
    "karma-coverage": "~2.2.0",
    "karma-jasmine": "~5.1.0",
    "karma-jasmine-html-reporter": "~2.1.0",
    "typescript": "~5.2.0"
  }
}`, g.config.Name)

	return os.WriteFile(filepath.Join(g.config.Path, "web/package.json"), []byte(content), 0644)
}

func (g *Generator) generateAngularJSON() error {
	content := fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "newProjectRoot": "projects",
  "projects": {
    "%s": {
      "projectType": "application",
      "root": "",
      "sourceRoot": "src",
      "prefix": "app",
      "architect": {
        "build": {
          "builder": "@angular-devkit/build-angular:application",
          "options": {
            "outputPath": "dist",
            "index": "src/index.html",
            "browser": "src/main.ts",
            "polyfills": ["zone.js"],
            "tsConfig": "tsconfig.app.json",
            "assets": ["src/assets"],
            "styles": ["src/styles.css"],
            "scripts": []
          },
          "configurations": {
            "production": {
              "budgets": [
                {
                  "type": "initial",
                  "maximumWarning": "500kb",
                  "maximumError": "1mb"
                }
              ],
              "outputHashing": "all",
              "fileReplacements": [
                {
                  "replace": "src/environments/environment.ts",
                  "with": "src/environments/environment.prod.ts"
                }
              ]
            },
            "development": {
              "optimization": false,
              "extractLicenses": false,
              "sourceMap": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": {
              "buildTarget": "%s:build:production"
            },
            "development": {
              "buildTarget": "%s:build:development"
            }
          },
          "defaultConfiguration": "development"
        },
        "test": {
          "builder": "@angular-devkit/build-angular:karma",
          "options": {
            "polyfills": ["zone.js", "zone.js/testing"],
            "tsConfig": "tsconfig.spec.json",
            "assets": ["src/assets"],
            "styles": ["src/styles.css"],
            "scripts": []
          }
        }
      }
    }
  }
}`, g.config.Name, g.config.Name, g.config.Name)

	return os.WriteFile(filepath.Join(g.config.Path, "web/angular.json"), []byte(content), 0644)
}

func (g *Generator) generateTsConfig() error {
	content := `{
  "compileOnSave": false,
  "compilerOptions": {
    "outDir": "./dist/out-tsc",
    "strict": true,
    "noImplicitOverride": true,
    "noPropertyAccessFromIndexSignature": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "sourceMap": true,
    "declaration": false,
    "experimentalDecorators": true,
    "moduleResolution": "bundler",
    "importHelpers": true,
    "target": "ES2022",
    "module": "ES2022",
    "useDefineForClassFields": false,
    "lib": ["ES2022", "dom"],
    "baseUrl": "./src",
    "paths": {
      "@core/*": ["app/core/*"],
      "@shared/*": ["app/shared/*"],
      "@features/*": ["app/features/*"],
      "@env/*": ["environments/*"]
    }
  },
  "angularCompilerOptions": {
    "enableI18nLegacyMessageIdFormat": false,
    "strictInjectionParameters": true,
    "strictInputAccessModifiers": true,
    "strictTemplates": true
  }
}`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/tsconfig.json"), []byte(content), 0644); err != nil {
		return err
	}

	appConfig := `{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "outDir": "./out-tsc/app",
    "types": []
  },
  "files": ["src/main.ts"],
  "include": ["src/**/*.d.ts"]
}`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/tsconfig.app.json"), []byte(appConfig), 0644); err != nil {
		return err
	}

	specConfig := `{
  "extends": "./tsconfig.json",
  "compilerOptions": {
    "outDir": "./out-tsc/spec",
    "types": ["jasmine"]
  },
  "include": ["src/**/*.spec.ts", "src/**/*.d.ts"]
}`

	return os.WriteFile(filepath.Join(g.config.Path, "web/tsconfig.spec.json"), []byte(specConfig), 0644)
}

func (g *Generator) generateAppModule() error {
	/* Generate main.ts */
	mainTs := `/*
 * GoAstra Frontend - Application Bootstrap
 */
import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/main.ts"), []byte(mainTs), 0644); err != nil {
		return err
	}

	/* Generate app.component.ts */
	appComponent := `/*
 * GoAstra Frontend - Root Component
 */
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: ` + "`" + `
    <router-outlet></router-outlet>
  ` + "`" + `,
  styles: []
})
export class AppComponent {
  title = 'GoAstra';
}
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/app.component.ts"), []byte(appComponent), 0644); err != nil {
		return err
	}

	/* Generate app.config.ts */
	appConfig := `/*
 * GoAstra Frontend - Application Configuration
 */
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { routes } from './app.routes';
import { authInterceptor } from '@core/interceptors/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor]))
  ]
};
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/app.config.ts"), []byte(appConfig), 0644); err != nil {
		return err
	}

	/* Generate app.routes.ts */
	appRoutes := `/*
 * GoAstra Frontend - Route Configuration
 */
import { Routes } from '@angular/router';
import { authGuard } from '@core/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full'
  },
  {
    path: 'home',
    loadComponent: () => import('@features/home/home.component').then(m => m.HomeComponent)
  },
  {
    path: 'auth',
    loadChildren: () => import('@features/auth/auth.routes').then(m => m.AUTH_ROUTES)
  },
  {
    path: '**',
    redirectTo: 'home'
  }
];
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/app.routes.ts"), []byte(appRoutes), 0644); err != nil {
		return err
	}

	/* Generate index.html */
	indexHtml := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <base href="/">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
</head>
<body>
  <app-root></app-root>
</body>
</html>
`, g.config.Name)

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/index.html"), []byte(indexHtml), 0644); err != nil {
		return err
	}

	/* Generate styles.css */
	styles := `/* GoAstra Frontend - Global Styles */

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  line-height: 1.5;
  color: #333;
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "web/src/styles.css"), []byte(styles), 0644)
}

func (g *Generator) generateCoreModule() error {
	/* Generate auth guard */
	authGuard := `/*
 * GoAstra Frontend - Auth Guard
 *
 * Protects routes requiring authentication.
 */
import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthService } from '@core/services/auth.service';

export const authGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (authService.isAuthenticated()) {
    return true;
  }

  router.navigate(['/auth/login'], {
    queryParams: { returnUrl: state.url }
  });

  return false;
};
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/core/guards/auth.guard.ts"), []byte(authGuard), 0644); err != nil {
		return err
	}

	/* Generate auth interceptor */
	authInterceptor := `/*
 * GoAstra Frontend - Auth Interceptor
 *
 * Attaches JWT token to outgoing HTTP requests.
 */
import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '@core/services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.getToken();

  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: ` + "`Bearer ${token}`" + `
      }
    });
  }

  return next(req);
};
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/core/interceptors/auth.interceptor.ts"), []byte(authInterceptor), 0644); err != nil {
		return err
	}

	/* Generate auth service */
	authService := `/*
 * GoAstra Frontend - Auth Service
 *
 * Handles authentication state and API calls.
 */
import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { environment } from '@env/environment';

interface AuthResponse {
  token: string;
  refreshToken: string;
  user: User;
}

interface User {
  id: number;
  email: string;
  name: string;
  role: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly TOKEN_KEY = 'auth_token';
  private readonly USER_KEY = 'auth_user';

  currentUser = signal<User | null>(null);

  constructor(
    private http: HttpClient,
    private router: Router
  ) {
    this.loadStoredUser();
  }

  login(email: string, password: string): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(` + "`${environment.apiUrl}/auth/login`" + `, { email, password })
      .pipe(tap(response => this.handleAuthResponse(response)));
  }

  register(data: { email: string; password: string; name: string }): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(` + "`${environment.apiUrl}/auth/register`" + `, data)
      .pipe(tap(response => this.handleAuthResponse(response)));
  }

  logout(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
    this.currentUser.set(null);
    this.router.navigate(['/auth/login']);
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  private handleAuthResponse(response: AuthResponse): void {
    localStorage.setItem(this.TOKEN_KEY, response.token);
    localStorage.setItem(this.USER_KEY, JSON.stringify(response.user));
    this.currentUser.set(response.user);
  }

  private loadStoredUser(): void {
    const stored = localStorage.getItem(this.USER_KEY);
    if (stored) {
      this.currentUser.set(JSON.parse(stored));
    }
  }
}
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/app/core/services/auth.service.ts"), []byte(authService), 0644); err != nil {
		return err
	}

	/* Generate API service */
	apiService := `/*
 * GoAstra Frontend - API Service
 *
 * Base service for HTTP API calls with error handling.
 */
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from '@env/environment';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get<T>(path: string, params?: Record<string, string>): Observable<T> {
    let httpParams = new HttpParams();
    if (params) {
      Object.keys(params).forEach(key => {
        httpParams = httpParams.set(key, params[key]);
      });
    }

    return this.http.get<T>(` + "`${this.baseUrl}${path}`" + `, { params: httpParams })
      .pipe(catchError(this.handleError));
  }

  post<T>(path: string, body: unknown): Observable<T> {
    return this.http.post<T>(` + "`${this.baseUrl}${path}`" + `, body)
      .pipe(catchError(this.handleError));
  }

  put<T>(path: string, body: unknown): Observable<T> {
    return this.http.put<T>(` + "`${this.baseUrl}${path}`" + `, body)
      .pipe(catchError(this.handleError));
  }

  delete<T>(path: string): Observable<T> {
    return this.http.delete<T>(` + "`${this.baseUrl}${path}`" + `)
      .pipe(catchError(this.handleError));
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'An error occurred';

    if (error.error instanceof ErrorEvent) {
      errorMessage = error.error.message;
    } else {
      errorMessage = error.error?.message || ` + "`Error: ${error.status}`" + `;
    }

    console.error(errorMessage);
    return throwError(() => new Error(errorMessage));
  }
}
`

	return os.WriteFile(filepath.Join(g.config.Path, "web/src/app/core/services/api.service.ts"), []byte(apiService), 0644)
}

func (g *Generator) generateEnvironments() error {
	envDev := `/*
 * GoAstra Frontend - Development Environment
 */
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1'
};
`

	if err := os.WriteFile(filepath.Join(g.config.Path, "web/src/environments/environment.ts"), []byte(envDev), 0644); err != nil {
		return err
	}

	envProd := `/*
 * GoAstra Frontend - Production Environment
 */
export const environment = {
  production: true,
  apiUrl: '/api/v1'
};
`

	return os.WriteFile(filepath.Join(g.config.Path, "web/src/environments/environment.prod.ts"), []byte(envProd), 0644)
}
