/***
 * GoAstra CLI - Migration Types and Configuration
 *
 * Defines the core data structures, configuration types, and driver
 * constants used throughout the GoAstra migration system.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package migrator

import (
	"database/sql"
	"time"
)

/***
 * Supported database driver constants.
 * GoAstra supports PostgreSQL, MySQL/MariaDB, and SQLite.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
const (
	DriverPostgres = "postgres"
	DriverMySQL    = "mysql"
	DriverSQLite   = "sqlite3"
)

/***
 * Migration represents a single database migration unit.
 * Encapsulates both the forward (up) and reverse (down) SQL operations
 * required for complete schema version control.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
type Migration struct {
	Version   string
	Name      string
	Filename  string
	Batch     int
	AppliedAt time.Time
}

/***
 * MigrationStatus provides a comprehensive view of a migration's state.
 * Used by the status command to display migration health and pending work.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
type MigrationStatus struct {
	Migration Migration
	Pending   bool
	Ran       bool
}

/***
 * Migrator is GoAstra's core migration orchestrator.
 * Coordinates all migration operations including discovery, execution,
 * rollback, and state tracking across the application lifecycle.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
type Migrator struct {
	db             *sql.DB
	migrationsPath string
	tableName      string
	driver         string
}

/***
 * Config defines the migration system configuration.
 * Allows customization of paths, database connections, and behavior.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
type Config struct {
	DatabaseURL    string
	MigrationsPath string
	TableName      string
	Driver         string
}

/***
 * DefaultConfig returns GoAstra's recommended migration configuration.
 * Follows the standard GoAstra project structure conventions.
 * Default driver is MySQL for localhost development.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func DefaultConfig() *Config {
	return &Config{
		MigrationsPath: "app/database/migrations",
		TableName:      "goastra_migrations",
		Driver:         DriverMySQL,
	}
}

/***
 * GetDriver returns the current database driver.
 * Useful for driver-specific SQL generation.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) GetDriver() string {
	return m.driver
}
