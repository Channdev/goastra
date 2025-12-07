package backend

func ConfigGo() string {
	return `package config

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
}

func LoggerGo() string {
	return `package logger

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
}

func DatabaseGo(db string) string {
	if db == "mysql" {
		return `package database

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
	}

	return `package database

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

func AuthGo() string {
	return `package auth

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
}

func MiddlewareGo() string {
	return `package middleware

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
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
`
}

func ModelsGo() string {
	return `package models

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
}

func SchemaTypesGo() string {
	return `package types

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
}

func SchemaGoMod() string {
	return `module schema

go 1.21
`
}

func HandlersGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is running",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
`
}

func RepositoryGo() string {
	return `package repository

import (
	"context"
)

type Repository[T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

type BaseRepository struct{}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}
`
}

func ServicesGo() string {
	return `package services

import (
	"context"
)

type Service[T any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

type BaseService struct{}

func NewBaseService() *BaseService {
	return &BaseService{}
}
`
}

func RouterGo() string {
	return `package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func New() *Router {
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{engine: engine}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) RegisterHealthRoutes(basePath string, handler interface {
	Health(c *gin.Context)
	Ready(c *gin.Context)
}) {
	health := r.engine.Group(basePath)
	{
		health.GET("/health", handler.Health)
		health.GET("/ready", handler.Ready)
	}
}

func (r *Router) RegisterAPIRoutes(basePath string, setupFunc func(rg *gin.RouterGroup)) {
	api := r.engine.Group(basePath)
	setupFunc(api)
}
`
}

func ValidatorGo() string {
	return `package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidations(v)
	}
}

func registerCustomValidations(v *validator.Validate) {
	// Add custom validations here
	// Example:
	// v.RegisterValidation("customtag", customValidationFunc)
}
`
}
