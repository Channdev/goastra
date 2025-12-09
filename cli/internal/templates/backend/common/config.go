/*
 * GoAstra CLI - Config Template
 *
 * Generates the application configuration loader.
 * Supports environment-based configuration with sensible defaults.
 */
package common

// ConfigGo returns the config.go template for loading app configuration.
func ConfigGo() string {
	return `package config

import (
	"fmt"
	"os"
)

/*
 * Config holds all application configuration values.
 * Values are loaded from environment variables with fallbacks.
 */
type Config struct {
	Env         string
	Port        string
	DBURL       string
	JWTSecret   string
	JWTExpiry   string
	CORSOrigins string
	LogLevel    string
}

/*
 * Load reads configuration from environment variables.
 * Returns a fully populated Config struct with defaults applied.
 */
func Load() *Config {
	cfg := &Config{
		Env:         getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DBURL:       getDatabaseURL(),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production-32chars"),
		JWTExpiry:   getEnv("JWT_EXPIRY", "24h"),
		CORSOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}
	return cfg
}

/*
 * IsProduction returns true if running in production environment.
 */
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

/*
 * IsDevelopment returns true if running in development environment.
 */
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

/*
 * getDatabaseURL builds database URL from environment variables.
 * Supports both DB_URL and individual MySQL/Postgres variables.
 */
func getDatabaseURL() string {
	// Check DB_URL first
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

	// Try MySQL individual vars
	if url := buildMySQLURL(); url != "" {
		return url
	}

	// Try Postgres individual vars
	if url := buildPostgresURL(); url != "" {
		return url
	}

	return ""
}

func buildMySQLURL() string {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	if host == "" || user == "" {
		return ""
	}
	if port == "" {
		port = "3306"
	}
	if db == "" {
		db = "goastra"
	}
	if pass != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, db)
	}
	return fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true", user, host, port, db)
}

func buildPostgresURL() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" {
		return ""
	}
	if port == "" {
		port = "5432"
	}
	if db == "" {
		db = "goastra"
	}
	if pass != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
	}
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", user, host, port, db)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
`
}
