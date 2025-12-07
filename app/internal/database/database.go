/*
 * GoAstra Backend - Database Connection
 *
 * Handles PostgreSQL connection pooling and lifecycle management.
 * Provides centralized database access for repositories.
 */
package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
 * DB wraps sqlx.DB with additional configuration.
 */
type DB struct {
	*sqlx.DB
}

/*
 * Connect establishes a connection to PostgreSQL.
 * Configures connection pool settings for production use.
 */
func Connect(databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return &DB{}, nil
	}

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	/* Configure connection pool */
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	/* Verify connection */
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

/*
 * ConnectWithConfig establishes connection using explicit configuration.
 */
func ConnectWithConfig(url string, maxOpen, maxIdle int, maxLife time.Duration) (*DB, error) {
	if url == "" {
		return &DB{}, nil
	}

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLife)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

/*
 * Close gracefully closes the database connection pool.
 */
func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}
	return db.DB.Close()
}

/*
 * Health checks if the database connection is alive.
 */
func (db *DB) Health() error {
	if db.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return db.Ping()
}

/*
 * Stats returns database connection pool statistics.
 */
func (db *DB) Stats() map[string]interface{} {
	if db.DB == nil {
		return nil
	}

	stats := db.DB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}
