/*
 * GoAstra CLI - Environment Loader
 *
 * Handles loading and parsing of environment-specific configuration files.
 * Supports .env.development, .env.production, and .env.test files.
 */
package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

/*
 * Environment represents the application runtime environment.
 * Used to determine which configuration file to load.
 */
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Test        Environment = "test"
)

/*
 * Config holds the parsed environment configuration.
 * Maps directly to the .env file structure.
 */
type Config struct {
	AppEnv    string
	APIURL    string
	DBURL     string
	JWTSecret string
	Port      string
	LogLevel  string
}

/*
 * Load reads the environment file for the specified environment.
 * Falls back to .env if specific file not found.
 */
func Load(envName string) error {
	envFile := fmt.Sprintf(".env.%s", envName)

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		if _, err := os.Stat(".env"); err == nil {
			envFile = ".env"
		} else {
			return nil
		}
	}

	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("failed to load %s: %w", envFile, err)
	}

	return nil
}

/*
 * LoadFromPath loads environment from a specific file path.
 * Used when running from non-project directories.
 */
func LoadFromPath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("environment file not found: %s", absPath)
	}

	return godotenv.Load(absPath)
}

/*
 * GetConfig returns the current environment configuration.
 * Reads from os.Environ after loading env file.
 */
func GetConfig() *Config {
	return &Config{
		AppEnv:    getEnvOrDefault("APP_ENV", "development"),
		APIURL:    getEnvOrDefault("API_URL", "http://localhost:8080"),
		DBURL:     getEnvOrDefault("DB_URL", ""),
		JWTSecret: getEnvOrDefault("JWT_SECRET", ""),
		Port:      getEnvOrDefault("PORT", "8080"),
		LogLevel:  getEnvOrDefault("LOG_LEVEL", "info"),
	}
}

/*
 * getEnvOrDefault retrieves an environment variable with fallback.
 */
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

/*
 * Validate checks that required environment variables are set.
 * Returns error with list of missing variables.
 */
func Validate(required []string) error {
	missing := []string{}

	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

/*
 * WriteEnvFile creates an environment file with the specified variables.
 * Used during project scaffolding.
 */
func WriteEnvFile(path string, vars map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create env file: %w", err)
	}
	defer file.Close()

	for key, value := range vars {
		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return fmt.Errorf("failed to write env variable: %w", err)
		}
	}

	return nil
}
