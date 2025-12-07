/*
 * GoAstra Backend - Configuration
 *
 * Centralized configuration management loaded from environment variables.
 * Provides type-safe access to all application settings.
 */
package config

import (
	"os"
	"strconv"
	"time"
)

/*
 * Config holds all application configuration values.
 * Values are loaded from environment variables with sensible defaults.
 */
type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	JWTSecret   string
	LogLevel    string
	APIURL      string

	/* JWT Configuration */
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration

	/* Database Pool Configuration */
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxLife  time.Duration

	/* CORS Configuration */
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
}

/*
 * Load reads configuration from environment variables.
 * Falls back to sensible defaults for development.
 */
func Load() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DB_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		APIURL:      getEnv("API_URL", "http://localhost:8080"),

		JWTExpiry:        getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		JWTRefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 7*24*time.Hour),

		DBMaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLife:  getDurationEnv("DB_CONN_MAX_LIFE", 5*time.Minute),

		CORSAllowedOrigins: getSliceEnv("CORS_ALLOWED_ORIGINS", []string{"*"}),
		CORSAllowedMethods: getSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		CORSAllowedHeaders: getSliceEnv("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Requested-With"}),
	}
}

/*
 * IsDevelopment returns true if running in development mode.
 */
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

/*
 * IsProduction returns true if running in production mode.
 */
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

/*
 * IsTest returns true if running in test mode.
 */
func (c *Config) IsTest() bool {
	return c.AppEnv == "test"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func getSliceEnv(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		/* Simple comma-separated parsing */
		result := []string{}
		current := ""
		for _, c := range value {
			if c == ',' {
				if current != "" {
					result = append(result, current)
					current = ""
				}
			} else {
				current += string(c)
			}
		}
		if current != "" {
			result = append(result, current)
		}
		return result
	}
	return fallback
}
