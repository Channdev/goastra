package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	skipAngular  bool
	skipBackend  bool
	useGraphQL   bool
	templateName string
	dbDriver     string
)

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new GoAstra project",
	Long: `Creates a new GoAstra project with Go backend and Angular frontend.

Templates:
  default   - Full-featured template with auth, dashboard, and beautiful landing page
  minimal   - Minimal starter template with basic structure

Database:
  postgres  - PostgreSQL (default)
  mysql     - MySQL/MariaDB`,
	Args: cobra.ExactArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVar(&skipAngular, "skip-angular", false, "skip Angular frontend generation")
	newCmd.Flags().BoolVar(&skipBackend, "skip-backend", false, "skip Go backend generation")
	newCmd.Flags().BoolVar(&useGraphQL, "graphql", false, "use GraphQL instead of REST")
	newCmd.Flags().StringVarP(&templateName, "template", "t", "default", "project template (default, minimal)")
	newCmd.Flags().StringVar(&dbDriver, "db", "postgres", "database driver (postgres, mysql)")
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	if err := validateProjectName(projectName); err != nil {
		return err
	}

	if templateName != "default" && templateName != "minimal" {
		return fmt.Errorf("invalid template: %s (use 'default' or 'minimal')", templateName)
	}

	if dbDriver != "postgres" && dbDriver != "mysql" {
		return fmt.Errorf("invalid database driver: %s (use 'postgres' or 'mysql')", dbDriver)
	}

	projectPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists: %s", projectPath)
	}

	color.Cyan("Creating new GoAstra project: %s\n", projectName)
	color.Cyan("Template: %s | Database: %s\n\n", templateName, dbDriver)

	color.Yellow("[1/7] Creating project structure...\n")
	if err := createDirectories(projectPath); err != nil {
		return err
	}

	color.Yellow("[2/7] Generating configuration files...\n")
	if err := generateConfigFiles(projectPath, projectName); err != nil {
		return err
	}

	color.Yellow("[3/7] Generating environment files...\n")
	if err := generateEnvFilesWithDB(projectPath, dbDriver); err != nil {
		return err
	}

	if !skipBackend {
		color.Yellow("[4/7] Generating Go backend...\n")
		if err := generateBackendWithDB(projectPath, projectName, dbDriver); err != nil {
			return err
		}
	}

	if !skipAngular {
		color.Yellow("[5/7] Generating Angular frontend...\n")
		if err := generateFrontendWithTemplate(projectPath, projectName, templateName); err != nil {
			return err
		}
	}

	color.Yellow("[6/7] Generating schema types...\n")
	if err := generateSchema(projectPath); err != nil {
		return err
	}

	color.Yellow("[7/7] Installing dependencies...\n")
	if err := installDependencies(projectPath, skipBackend, skipAngular); err != nil {
		color.Yellow("Warning: Failed to install some dependencies: %v\n", err)
		color.Yellow("You may need to run 'go mod tidy' in app/ and 'npm install' in web/ manually.\n")
	}

	color.Green("\nProject created successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Printf("  goastra dev\n\n")
	fmt.Printf("Your app will be available at:\n")
	fmt.Printf("  Frontend: http://localhost:4200\n")
	fmt.Printf("  Backend:  http://localhost:8080\n")

	return nil
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("project name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("project name too long (max 100 characters)")
	}
	for _, c := range name {
		if !isValidNameChar(c) {
			return fmt.Errorf("invalid character in project name: %c", c)
		}
	}
	return nil
}

func isValidNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}

func createDirectories(projectPath string) error {
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
		"web/src/app/features/home",
		"web/src/app/features/auth/login",
		"web/src/app/features/auth/register",
		"web/src/app/features/dashboard",
		"web/src/app/features/not-found",
		"web/src/environments",
		"web/src/assets",
		"schema/types",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func generateConfigFiles(projectPath, projectName string) error {
	goastraJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "api": {
    "type": "rest",
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
    "outputPath": "web/src/app/core/models"
  },
  "database": {
    "driver": "postgres",
    "migrationsPath": "app/migrations"
  }
}`, projectName, projectName)

	gitignore := `dist/
bin/
node_modules/
vendor/
.env
.env.local
.env.*.local
.idea/
.vscode/
*.exe
*.dll
*.so
*.dylib
web/dist/
web/.angular/
coverage/
*.log
tmp/
`

	if err := os.WriteFile(filepath.Join(projectPath, "goastra.json"), []byte(goastraJSON), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignore), 0644)
}

func generateEnvFiles(projectPath string) error {
	return generateEnvFilesWithDB(projectPath, "postgres")
}

func generateEnvFilesWithDB(projectPath, db string) error {
	var dbURLDev, dbURLTest string
	if db == "mysql" {
		dbURLDev = "user:password@tcp(localhost:3306)/goastra_dev?parseTime=true"
		dbURLTest = "user:password@tcp(localhost:3306)/goastra_test?parseTime=true"
	} else {
		dbURLDev = "postgres://user:password@localhost:5432/goastra_dev?sslmode=disable"
		dbURLTest = "postgres://user:password@localhost:5432/goastra_test?sslmode=disable"
	}

	envDev := fmt.Sprintf(`APP_ENV=development
API_URL=http://localhost:8080
PORT=8080
LOG_LEVEL=debug
DB_DRIVER=%s
DB_URL=%s
JWT_SECRET=dev-secret-change-in-production-32chars
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=http://localhost:4200
`, db, dbURLDev)

	envProd := fmt.Sprintf(`APP_ENV=production
API_URL=https://api.example.com
PORT=8080
LOG_LEVEL=info
DB_DRIVER=%s
DB_URL=
JWT_SECRET=
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://example.com
`, db)

	envTest := fmt.Sprintf(`APP_ENV=test
API_URL=http://localhost:8081
PORT=8081
LOG_LEVEL=error
DB_DRIVER=%s
DB_URL=%s
JWT_SECRET=test-secret-32-characters-long!!
JWT_EXPIRY=1h
CORS_ALLOWED_ORIGINS=*
`, db, dbURLTest)

	if err := os.WriteFile(filepath.Join(projectPath, ".env.development"), []byte(envDev), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(projectPath, ".env.production"), []byte(envProd), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, ".env.test"), []byte(envTest), 0644)
}

func generateBackend(projectPath, projectName string) error {
	goMod := fmt.Sprintf(`module github.com/%s/app

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.16.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.16.0
)
`, projectName)

	mainGo := `package main

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

	if os.Getenv("APP_ENV") == "production" {
		r.Static("/assets", "./public/browser/assets")
		r.StaticFile("/favicon.ico", "./public/browser/favicon.ico")
		r.NoRoute(func(c *gin.Context) {
			c.File("./public/browser/index.html")
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

	if err := os.WriteFile(filepath.Join(projectPath, "app/go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(projectPath, "app/cmd/server/main.go"), []byte(mainGo), 0644); err != nil {
		return err
	}

	return generateInternalFilesWithDB(projectPath, projectName, "postgres")
}

func generateBackendWithDB(projectPath, projectName, db string) error {
	dbImport := "github.com/lib/pq v1.10.9"
	if db == "mysql" {
		dbImport = "github.com/go-sql-driver/mysql v1.7.1"
	}

	goMod := fmt.Sprintf(`module github.com/%s/app

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

	mainGo := `package main

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
		r.Static("/assets", "./public/browser/assets")
		r.StaticFile("/favicon.ico", "./public/browser/favicon.ico")
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			indexPath := filepath.Join(".", "public", "browser", "index.html")
			c.File(indexPath)
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

	if err := os.WriteFile(filepath.Join(projectPath, "app/go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(projectPath, "app/cmd/server/main.go"), []byte(mainGo), 0644); err != nil {
		return err
	}

	return generateInternalFilesWithDB(projectPath, projectName, db)
}

func generateInternalFiles(projectPath, projectName string) error {
	return generateInternalFilesWithDB(projectPath, projectName, "postgres")
}

func generateInternalFilesWithDB(projectPath, projectName, db string) error {
	configGo := `package config

import "os"

type Config struct {
	Env         string
	Port        string
	DBDriver    string
	DBURL       string
	JWTSecret   string
	JWTExpiry   string
	CORSOrigins string
}

func Load() *Config {
	return &Config{
		Env:         getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DBDriver:    getEnv("DB_DRIVER", "postgres"),
		DBURL:       getEnv("DB_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry:   getEnv("JWT_EXPIRY", "24h"),
		CORSOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
`

	loggerGo := `package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(env string) *Logger {
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, _ := config.Build()
	return &Logger{logger.Sugar()}
}
`

	var databaseGo string
	if db == "mysql" {
		databaseGo = `package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sqlx.DB
}

func Connect(url string) (*DB, error) {
	if url == "" {
		return nil, nil
	}
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}
	db, err := sqlx.Connect(driver, url)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Health() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.Ping()
}

func (db *DB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}
`
	} else {
		databaseGo = `package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(url string) (*DB, error) {
	if url == "" {
		return nil, nil
	}
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}
	db, err := sqlx.Connect(driver, url)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Health() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.Ping()
}

func (db *DB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}
`
	}

	authGo := `package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint   ` + "`json:\"user_id\"`" + `
	Email  string ` + "`json:\"email\"`" + `
	Role   string ` + "`json:\"role\"`" + `
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, role, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
`

	middlewareGo := `package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins == "*" || strings.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		c.Set("token", parts[1])
		c.Next()
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateID()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func generateID() string {
	return "req_" + randomString(16)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%%len(letters)]
	}
	return string(b)
}
`

	modelsGo := `package models

import "time"

type User struct {
	ID        uint      ` + "`json:\"id\" db:\"id\"`" + `
	Email     string    ` + "`json:\"email\" db:\"email\"`" + `
	Password  string    ` + "`json:\"-\" db:\"password\"`" + `
	Name      string    ` + "`json:\"name\" db:\"name\"`" + `
	Role      string    ` + "`json:\"role\" db:\"role\"`" + `
	Active    bool      ` + "`json:\"active\" db:\"active\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
}

type LoginRequest struct {
	Email    string ` + "`json:\"email\" binding:\"required,email\"`" + `
	Password string ` + "`json:\"password\" binding:\"required,min=6\"`" + `
}

type RegisterRequest struct {
	Email    string ` + "`json:\"email\" binding:\"required,email\"`" + `
	Password string ` + "`json:\"password\" binding:\"required,min=6\"`" + `
	Name     string ` + "`json:\"name\" binding:\"required\"`" + `
}

type AuthResponse struct {
	Token     string ` + "`json:\"token\"`" + `
	ExpiresAt int64  ` + "`json:\"expires_at\"`" + `
	User      *User  ` + "`json:\"user\"`" + `
}
`

	files := map[string]string{
		"app/internal/config/config.go":         configGo,
		"app/internal/logger/logger.go":         loggerGo,
		"app/internal/database/database.go":     databaseGo,
		"app/internal/auth/auth.go":             authGo,
		"app/internal/middleware/middleware.go": middlewareGo,
		"app/internal/models/models.go":         modelsGo,
	}

	for path, content := range files {
		fullPath := filepath.Join(projectPath, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func generateFrontendWithTemplate(projectPath, projectName, template string) error {
	if template == "minimal" {
		return generateFrontendMinimal(projectPath, projectName)
	}
	return generateFrontendDefault(projectPath, projectName)
}

func generateFrontendMinimal(projectPath, projectName string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "test": "ng test"
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
    "typescript": "~5.2.0"
  }
}`, projectName)

	angularJSON := fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "cli": { "analytics": false },
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
            "styles": ["src/styles.css"]
          },
          "configurations": {
            "production": {
              "outputHashing": "all",
              "fileReplacements": [{
                "replace": "src/environments/environment.ts",
                "with": "src/environments/environment.prod.ts"
              }]
            },
            "development": {
              "optimization": false,
              "sourceMap": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": { "buildTarget": "%s:build:production" },
            "development": { "buildTarget": "%s:build:development" }
          },
          "defaultConfiguration": "development"
        }
      }
    }
  }
}`, projectName, projectName, projectName)

	tsconfig := `{
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
    "lib": ["ES2022", "dom"]
  },
  "angularCompilerOptions": {
    "enableI18nLegacyMessageIdFormat": false,
    "strictInjectionParameters": true,
    "strictInputAccessModifiers": true,
    "strictTemplates": true
  }
}`

	tsconfigApp := `{
  "extends": "./tsconfig.json",
  "compilerOptions": { "outDir": "./out-tsc/app" },
  "files": ["src/main.ts"],
  "include": ["src/**/*.d.ts"]
}`

	proxyConf := `{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false,
    "changeOrigin": true
  }
}`

	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <base href="/">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
  <app-root></app-root>
</body>
</html>`, projectName)

	mainTS := `import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig).catch((err) => console.error(err));
`

	stylesCSS := `* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f8fafc;
  color: #1e293b;
}
`

	appComponent := `import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: '<router-outlet></router-outlet>'
})
export class AppComponent {}
`

	appConfig := `import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [provideRouter(routes), provideHttpClient()]
};
`

	appRoutes := `import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: '**', redirectTo: '' }
];
`

	homeComponent := `import { Component } from '@angular/core';

@Component({
  selector: 'app-home',
  standalone: true,
  template: ` + "`" + `
    <main>
      <h1>Hello, GoAstra!</h1>
      <p>Start building your app.</p>
    </main>
  ` + "`" + `,
  styles: [` + "`" + `
    main { padding: 2rem; text-align: center; }
    h1 { margin-bottom: 0.5rem; }
  ` + "`" + `]
})
export class HomeComponent {}
`

	envDev := `export const environment = { production: false, apiUrl: 'http://localhost:8080/api/v1' };`
	envProd := `export const environment = { production: true, apiUrl: '/api/v1' };`

	if err := os.MkdirAll(filepath.Join(projectPath, "web/src/app/home"), 0755); err != nil {
		return err
	}

	files := map[string]string{
		"web/package.json":                         packageJSON,
		"web/angular.json":                         angularJSON,
		"web/tsconfig.json":                        tsconfig,
		"web/tsconfig.app.json":                    tsconfigApp,
		"web/proxy.conf.json":                      proxyConf,
		"web/src/index.html":                       indexHTML,
		"web/src/main.ts":                          mainTS,
		"web/src/styles.css":                       stylesCSS,
		"web/src/app/app.component.ts":             appComponent,
		"web/src/app/app.config.ts":                appConfig,
		"web/src/app/app.routes.ts":                appRoutes,
		"web/src/app/home/home.component.ts":       homeComponent,
		"web/src/environments/environment.ts":      envDev,
		"web/src/environments/environment.prod.ts": envProd,
	}

	for path, content := range files {
		fullPath := filepath.Join(projectPath, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func generateFrontendDefault(projectPath, projectName string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s-web",
  "version": "1.0.0",
  "scripts": {
    "ng": "ng",
    "start": "ng serve --proxy-config proxy.conf.json",
    "build": "ng build",
    "test": "ng test"
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
    "typescript": "~5.2.0"
  }
}`, projectName)

	angularJSON := fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "cli": { "analytics": false },
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
            "styles": ["src/styles.css"]
          },
          "configurations": {
            "production": {
              "outputHashing": "all",
              "fileReplacements": [{
                "replace": "src/environments/environment.ts",
                "with": "src/environments/environment.prod.ts"
              }]
            },
            "development": {
              "optimization": false,
              "sourceMap": true
            }
          },
          "defaultConfiguration": "production"
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server",
          "configurations": {
            "production": { "buildTarget": "%s:build:production" },
            "development": { "buildTarget": "%s:build:development" }
          },
          "defaultConfiguration": "development"
        }
      }
    }
  }
}`, projectName, projectName, projectName)

	tsconfig := `{
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

	tsconfigApp := `{
  "extends": "./tsconfig.json",
  "compilerOptions": { "outDir": "./out-tsc/app" },
  "files": ["src/main.ts"],
  "include": ["src/**/*.d.ts"]
}`

	proxyConf := `{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false,
    "changeOrigin": true
  }
}`

	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <base href="/">
  <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
  <app-root></app-root>
</body>
</html>`, projectName)

	mainTS := `import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

bootstrapApplication(AppComponent, appConfig).catch((err) => console.error(err));
`

	stylesCSS := `:root {
  --color-primary: #3b82f6;
  --color-background: #0f172a;
  --color-surface: #1e293b;
  --color-text: #f8fafc;
  --color-text-muted: #94a3b8;
  --color-border: #334155;
}

* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--color-background);
  color: var(--color-text);
}

a { color: var(--color-primary); text-decoration: none; }
`

	appComponent := `import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet],
  template: '<router-outlet></router-outlet>'
})
export class AppComponent {}
`

	appConfig := `import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [provideRouter(routes), provideHttpClient()]
};
`

	appRoutes := `import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', loadComponent: () => import('@features/home/home.component').then(m => m.HomeComponent) },
  { path: 'login', loadComponent: () => import('@features/auth/login/login.component').then(m => m.LoginComponent) },
  { path: 'register', loadComponent: () => import('@features/auth/register/register.component').then(m => m.RegisterComponent) },
  { path: 'dashboard', loadComponent: () => import('@features/dashboard/dashboard.component').then(m => m.DashboardComponent) },
  { path: '**', redirectTo: 'home' }
];
`

	homeComponent := `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterLink],
  template: ` + "`" + `
    <div class="landing">
      <nav class="navbar">
        <div class="logo">GoAstra</div>
        <div class="nav-links">
          <a routerLink="/login">Login</a>
          <a routerLink="/register" class="btn-primary">Get Started</a>
        </div>
      </nav>

      <main class="hero">
        <div class="hero-content">
          <h1>Build Full-Stack Apps<br><span class="gradient">Lightning Fast</span></h1>
          <p>GoAstra combines the power of Go backend with Angular frontend.<br>Production-ready, type-safe, and developer-friendly.</p>
          <div class="hero-buttons">
            <a routerLink="/register" class="btn btn-primary">Start Building</a>
            <a href="https://github.com/channdev/goastra" target="_blank" class="btn btn-secondary">View on GitHub</a>
          </div>
        </div>
      </main>

      <section class="features">
        <div class="feature">
          <div class="feature-icon">&#9889;</div>
          <h3>Blazing Fast</h3>
          <p>Go's performance meets Angular's reactivity for lightning-fast apps.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">&#128274;</div>
          <h3>Type Safe</h3>
          <p>End-to-end type safety with shared schemas between frontend and backend.</p>
        </div>
        <div class="feature">
          <div class="feature-icon">&#128640;</div>
          <h3>Production Ready</h3>
          <p>JWT auth, CORS, logging, and database support out of the box.</p>
        </div>
      </section>

      <footer class="footer">
        <p>Built with GoAstra &middot; <a href="https://github.com/channdev/goastra">GitHub</a></p>
      </footer>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .landing { min-height: 100vh; display: flex; flex-direction: column; }
    .navbar { display: flex; justify-content: space-between; align-items: center; padding: 1.5rem 3rem; }
    .logo { font-size: 1.5rem; font-weight: 700; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .nav-links { display: flex; gap: 1.5rem; align-items: center; }
    .nav-links a { color: #94a3b8; transition: color 0.2s; }
    .nav-links a:hover { color: #f8fafc; }
    .btn-primary { background: #3b82f6 !important; color: white !important; padding: 0.5rem 1rem; border-radius: 6px; }
    .hero { flex: 1; display: flex; align-items: center; justify-content: center; text-align: center; padding: 2rem; }
    .hero h1 { font-size: 3.5rem; line-height: 1.1; margin-bottom: 1.5rem; }
    .gradient { background: linear-gradient(135deg, #3b82f6, #8b5cf6, #ec4899); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .hero p { color: #94a3b8; font-size: 1.25rem; margin-bottom: 2rem; line-height: 1.6; }
    .hero-buttons { display: flex; gap: 1rem; justify-content: center; }
    .btn { padding: 0.875rem 1.75rem; border-radius: 8px; font-weight: 500; transition: transform 0.2s, box-shadow 0.2s; }
    .btn:hover { transform: translateY(-2px); }
    .btn-secondary { background: #1e293b; color: #f8fafc; border: 1px solid #334155; }
    .features { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 2rem; padding: 4rem 3rem; background: #1e293b; }
    .feature { text-align: center; padding: 2rem; }
    .feature-icon { font-size: 2.5rem; margin-bottom: 1rem; }
    .feature h3 { margin-bottom: 0.5rem; }
    .feature p { color: #94a3b8; }
    .footer { padding: 2rem; text-align: center; color: #64748b; border-top: 1px solid #334155; }
    .footer a { color: #3b82f6; }
  ` + "`" + `]
})
export class HomeComponent {}
`

	envDev := `export const environment = { production: false, apiUrl: 'http://localhost:8080/api/v1' };`
	envProd := `export const environment = { production: true, apiUrl: '/api/v1' };`

	loginComponent := `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterLink, FormsModule],
  template: ` + "`" + `
    <div class="auth-page">
      <div class="auth-card">
        <h1>Welcome Back</h1>
        <p class="subtitle">Sign in to your account</p>
        <form (ngSubmit)="onSubmit()">
          <div class="form-group">
            <label>Email</label>
            <input type="email" [(ngModel)]="email" name="email" placeholder="you@example.com" required>
          </div>
          <div class="form-group">
            <label>Password</label>
            <input type="password" [(ngModel)]="password" name="password" placeholder="Enter password" required>
          </div>
          <button type="submit" class="btn-submit">Sign In</button>
        </form>
        <p class="switch">Don't have an account? <a routerLink="/register">Sign up</a></p>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; }
    .auth-card { background: #1e293b; padding: 2.5rem; border-radius: 12px; width: 100%; max-width: 400px; }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: #94a3b8; margin-bottom: 2rem; }
    .form-group { margin-bottom: 1.25rem; }
    .form-group label { display: block; margin-bottom: 0.5rem; color: #94a3b8; font-size: 0.875rem; }
    .form-group input { width: 100%; padding: 0.75rem; background: #0f172a; border: 1px solid #334155; border-radius: 6px; color: #f8fafc; font-size: 1rem; }
    .form-group input:focus { outline: none; border-color: #3b82f6; }
    .btn-submit { width: 100%; padding: 0.875rem; background: #3b82f6; color: white; border: none; border-radius: 6px; font-size: 1rem; cursor: pointer; margin-top: 0.5rem; }
    .btn-submit:hover { background: #2563eb; }
    .switch { text-align: center; margin-top: 1.5rem; color: #94a3b8; }
    .switch a { color: #3b82f6; }
  ` + "`" + `]
})
export class LoginComponent {
  email = '';
  password = '';
  onSubmit() { console.log('Login:', this.email); }
}
`

	registerComponent := `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [RouterLink, FormsModule],
  template: ` + "`" + `
    <div class="auth-page">
      <div class="auth-card">
        <h1>Create Account</h1>
        <p class="subtitle">Start building with GoAstra</p>
        <form (ngSubmit)="onSubmit()">
          <div class="form-group">
            <label>Name</label>
            <input type="text" [(ngModel)]="name" name="name" placeholder="Your name" required>
          </div>
          <div class="form-group">
            <label>Email</label>
            <input type="email" [(ngModel)]="email" name="email" placeholder="you@example.com" required>
          </div>
          <div class="form-group">
            <label>Password</label>
            <input type="password" [(ngModel)]="password" name="password" placeholder="Create password" required>
          </div>
          <button type="submit" class="btn-submit">Create Account</button>
        </form>
        <p class="switch">Already have an account? <a routerLink="/login">Sign in</a></p>
      </div>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .auth-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 2rem; }
    .auth-card { background: #1e293b; padding: 2.5rem; border-radius: 12px; width: 100%; max-width: 400px; }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: #94a3b8; margin-bottom: 2rem; }
    .form-group { margin-bottom: 1.25rem; }
    .form-group label { display: block; margin-bottom: 0.5rem; color: #94a3b8; font-size: 0.875rem; }
    .form-group input { width: 100%; padding: 0.75rem; background: #0f172a; border: 1px solid #334155; border-radius: 6px; color: #f8fafc; font-size: 1rem; }
    .form-group input:focus { outline: none; border-color: #3b82f6; }
    .btn-submit { width: 100%; padding: 0.875rem; background: #3b82f6; color: white; border: none; border-radius: 6px; font-size: 1rem; cursor: pointer; margin-top: 0.5rem; }
    .btn-submit:hover { background: #2563eb; }
    .switch { text-align: center; margin-top: 1.5rem; color: #94a3b8; }
    .switch a { color: #3b82f6; }
  ` + "`" + `]
})
export class RegisterComponent {
  name = '';
  email = '';
  password = '';
  onSubmit() { console.log('Register:', this.email); }
}
`

	dashboardComponent := `import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink],
  template: ` + "`" + `
    <div class="dashboard">
      <nav class="sidebar">
        <div class="logo">GoAstra</div>
        <div class="nav-items">
          <a routerLink="/dashboard" class="active">Dashboard</a>
          <a routerLink="/dashboard">Users</a>
          <a routerLink="/dashboard">Settings</a>
        </div>
        <a routerLink="/home" class="logout">Logout</a>
      </nav>
      <main class="content">
        <header>
          <h1>Dashboard</h1>
          <p>Welcome back! Here's an overview of your application.</p>
        </header>
        <div class="stats">
          <div class="stat-card"><h3>1,234</h3><p>Total Users</p></div>
          <div class="stat-card"><h3>567</h3><p>Active Today</p></div>
          <div class="stat-card"><h3>89%</h3><p>Uptime</p></div>
          <div class="stat-card"><h3>12ms</h3><p>Avg Response</p></div>
        </div>
      </main>
    </div>
  ` + "`" + `,
  styles: [` + "`" + `
    .dashboard { display: flex; min-height: 100vh; }
    .sidebar { width: 240px; background: #1e293b; padding: 1.5rem; display: flex; flex-direction: column; }
    .logo { font-size: 1.25rem; font-weight: 700; margin-bottom: 2rem; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
    .nav-items { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
    .nav-items a { padding: 0.75rem 1rem; border-radius: 6px; color: #94a3b8; transition: all 0.2s; }
    .nav-items a:hover, .nav-items a.active { background: #334155; color: #f8fafc; }
    .logout { color: #94a3b8; padding: 0.75rem 1rem; }
    .content { flex: 1; padding: 2rem; }
    header { margin-bottom: 2rem; }
    header h1 { margin-bottom: 0.5rem; }
    header p { color: #94a3b8; }
    .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; }
    .stat-card { background: #1e293b; padding: 1.5rem; border-radius: 8px; }
    .stat-card h3 { font-size: 2rem; margin-bottom: 0.25rem; }
    .stat-card p { color: #94a3b8; }
  ` + "`" + `]
})
export class DashboardComponent {}
`

	files := map[string]string{
		"web/package.json":                                    packageJSON,
		"web/angular.json":                                    angularJSON,
		"web/tsconfig.json":                                   tsconfig,
		"web/tsconfig.app.json":                               tsconfigApp,
		"web/proxy.conf.json":                                 proxyConf,
		"web/src/index.html":                                  indexHTML,
		"web/src/main.ts":                                     mainTS,
		"web/src/styles.css":                                  stylesCSS,
		"web/src/app/app.component.ts":                        appComponent,
		"web/src/app/app.config.ts":                           appConfig,
		"web/src/app/app.routes.ts":                           appRoutes,
		"web/src/app/features/home/home.component.ts":         homeComponent,
		"web/src/app/features/auth/login/login.component.ts":  loginComponent,
		"web/src/app/features/auth/register/register.component.ts": registerComponent,
		"web/src/app/features/dashboard/dashboard.component.ts": dashboardComponent,
		"web/src/environments/environment.ts":                 envDev,
		"web/src/environments/environment.prod.ts":            envProd,
	}

	for path, content := range files {
		fullPath := filepath.Join(projectPath, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func generateSchema(projectPath string) error {
	schemaGo := `package types

import "time"

type BaseModel struct {
	ID        uint      ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

type User struct {
	BaseModel
	Email  string ` + "`json:\"email\"`" + `
	Name   string ` + "`json:\"name\"`" + `
	Role   string ` + "`json:\"role\"`" + `
	Active bool   ` + "`json:\"active\"`" + `
}

type APIError struct {
	Code    string ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}
`

	goMod := `module schema

go 1.21
`

	if err := os.WriteFile(filepath.Join(projectPath, "schema/types/types.go"), []byte(schemaGo), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectPath, "schema/go.mod"), []byte(goMod), 0644)
}

func runCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installDependencies(projectPath string, skipBackend, skipAngular bool) error {
	if !skipBackend {
		color.Blue("  Running 'go mod tidy' in app/...\n")
		appDir := filepath.Join(projectPath, "app")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = appDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go mod tidy failed: %w", err)
		}
	}

	if !skipAngular {
		color.Blue("  Running 'npm install' in web/...\n")
		webDir := filepath.Join(projectPath, "web")
		var cmd *exec.Cmd
		if os.PathSeparator == '\\' && os.PathListSeparator == ';' {
			cmd = exec.Command("cmd", "/c", "npm", "install")
		} else {
			cmd = exec.Command("npm", "install")
		}
		cmd.Dir = webDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}
	}

	return nil
}
