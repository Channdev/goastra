/***
 * GoAstra CLI - Database Migration Engine
 *
 * GoAstra's native database migration system built from the ground up for
 * Go applications. Provides enterprise-grade schema version control with
 * multi-driver support, transaction safety, and intelligent tracking.
 *
 * Supported Databases:
 *   - MySQL / MariaDB (default for localhost)
 *   - PostgreSQL
 *   - SQLite
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

/***
 * New instantiates a configured Migrator ready for operations.
 * Validates the configuration and initializes the migration environment.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func New(cfg *Config) (*Migrator, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return &Migrator{
		migrationsPath: cfg.MigrationsPath,
		tableName:      cfg.TableName,
		driver:         cfg.Driver,
	}, nil
}

/***
 * Connect establishes the database connection for migration operations.
 * This must be called before executing any migration commands.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Connect(databaseURL string) error {
	db, err := sql.Open(m.driver, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	m.db = db
	return nil
}

/***
 * Close terminates the database connection gracefully.
 * Should be called when migration operations are complete.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

/***
 * EnsureMigrationTable creates the GoAstra migrations tracking table.
 * This table maintains the history of all applied migrations.
 * Uses driver-specific SQL syntax for compatibility.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) EnsureMigrationTable() error {
	query := m.buildCreateTableQuery()

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

/***
 * buildCreateTableQuery generates driver-specific CREATE TABLE SQL.
 * Handles syntax differences between MySQL, PostgreSQL, and SQLite.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) buildCreateTableQuery() string {
	switch m.driver {
	case DriverMySQL:
		return fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id INT AUTO_INCREMENT PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				name VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`, m.tableName)

	case DriverPostgres:
		return fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				name VARCHAR(255) NOT NULL,
				batch INTEGER NOT NULL,
				applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			)
		`, m.tableName)

	case DriverSQLite:
		return fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				version VARCHAR(255) NOT NULL UNIQUE,
				name VARCHAR(255) NOT NULL,
				batch INTEGER NOT NULL,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)
		`, m.tableName)

	default:
		// Fallback to MySQL syntax
		return fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id INT AUTO_INCREMENT PRIMARY KEY,
				version VARCHAR(255) NOT NULL UNIQUE,
				name VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		`, m.tableName)
	}
}

/***
 * DiscoverMigrations scans the migrations directory for all migration files.
 * Supports both SQL and Go-based migrations with timestamp versioning.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) DiscoverMigrations() ([]Migration, error) {
	pattern := filepath.Join(m.migrationsPath, "*.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob migration files: %w", err)
	}

	goPattern := filepath.Join(m.migrationsPath, "*.go")
	goFiles, err := filepath.Glob(goPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob go migration files: %w", err)
	}

	files = append(files, goFiles...)

	var migrations []Migration
	migrationRegex := regexp.MustCompile(`^(\d{14})_(.+)\.(sql|go)$`)

	for _, file := range files {
		basename := filepath.Base(file)
		matches := migrationRegex.FindStringSubmatch(basename)
		if matches == nil {
			continue
		}

		migrations = append(migrations, Migration{
			Version:  matches[1],
			Name:     matches[2],
			Filename: file,
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

/***
 * GetAppliedMigrations retrieves all successfully executed migrations.
 * Returns migrations ordered by version for chronological tracking.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
	query := fmt.Sprintf(`
		SELECT version, name, batch, applied_at
		FROM %s
		ORDER BY version ASC
	`, m.tableName)

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var mig Migration
		if err := rows.Scan(&mig.Version, &mig.Name, &mig.Batch, &mig.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		migrations = append(migrations, mig)
	}

	return migrations, rows.Err()
}

/***
 * GetPendingMigrations identifies migrations awaiting execution.
 * Compares discovered migration files against the database history.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) GetPendingMigrations() ([]Migration, error) {
	allMigrations, err := m.DiscoverMigrations()
	if err != nil {
		return nil, err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, mig := range applied {
		appliedMap[mig.Version] = true
	}

	var pending []Migration
	for _, mig := range allMigrations {
		if !appliedMap[mig.Version] {
			pending = append(pending, mig)
		}
	}

	return pending, nil
}

/***
 * GetNextBatch calculates the next batch number for migration grouping.
 * Batch numbers enable GoAstra's grouped rollback functionality.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) GetNextBatch() (int, error) {
	query := fmt.Sprintf(`SELECT COALESCE(MAX(batch), 0) + 1 FROM %s`, m.tableName)

	var batch int
	err := m.db.QueryRow(query).Scan(&batch)
	if err != nil {
		return 0, fmt.Errorf("failed to get next batch number: %w", err)
	}

	return batch, nil
}

/***
 * GetLastBatch retrieves the most recent migration batch number.
 * Used by rollback operations to identify the target batch.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) GetLastBatch() (int, error) {
	query := fmt.Sprintf(`SELECT COALESCE(MAX(batch), 0) FROM %s`, m.tableName)

	var batch int
	err := m.db.QueryRow(query).Scan(&batch)
	if err != nil {
		return 0, fmt.Errorf("failed to get last batch number: %w", err)
	}

	return batch, nil
}

/***
 * runMigration executes a single migration in the specified direction.
 * Manages transaction boundaries and updates the migration registry.
 * Uses driver-specific placeholder syntax for SQL parameters.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) runMigration(mig Migration, batch int, direction string) error {
	content, err := os.ReadFile(mig.Filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	sqlContent := extractSQL(string(content), direction)
	if sqlContent == "" {
		return fmt.Errorf("no %s SQL found in migration", direction)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	if _, err := tx.Exec(sqlContent); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	if direction == "up" {
		if err := m.recordMigration(tx, mig, batch); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := m.removeMigrationRecord(tx, mig); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

/***
 * recordMigration inserts a migration record into the tracking table.
 * Uses driver-specific placeholder syntax.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) recordMigration(tx *sql.Tx, mig Migration, batch int) error {
	var query string

	switch m.driver {
	case DriverPostgres:
		query = fmt.Sprintf(
			`INSERT INTO %s (version, name, batch) VALUES ($1, $2, $3)`,
			m.tableName,
		)
	default:
		// MySQL, SQLite use ? placeholders
		query = fmt.Sprintf(
			`INSERT INTO %s (version, name, batch) VALUES (?, ?, ?)`,
			m.tableName,
		)
	}

	_, err := tx.Exec(query, mig.Version, mig.Name, batch)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

/***
 * removeMigrationRecord deletes a migration record from the tracking table.
 * Uses driver-specific placeholder syntax.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) removeMigrationRecord(tx *sql.Tx, mig Migration) error {
	var query string

	switch m.driver {
	case DriverPostgres:
		query = fmt.Sprintf(`DELETE FROM %s WHERE version = $1`, m.tableName)
	default:
		// MySQL, SQLite use ? placeholders
		query = fmt.Sprintf(`DELETE FROM %s WHERE version = ?`, m.tableName)
	}

	_, err := tx.Exec(query, mig.Version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return nil
}

/***
 * getMigrationFilename resolves the full filesystem path for a migration.
 * Searches for both SQL and Go-based migration file formats.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) getMigrationFilename(version, name string) string {
	sqlPath := filepath.Join(m.migrationsPath, fmt.Sprintf("%s_%s.sql", version, name))
	if _, err := os.Stat(sqlPath); err == nil {
		return sqlPath
	}

	goPath := filepath.Join(m.migrationsPath, fmt.Sprintf("%s_%s.go", version, name))
	if _, err := os.Stat(goPath); err == nil {
		return goPath
	}

	return sqlPath
}
