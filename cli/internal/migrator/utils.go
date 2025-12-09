/***
 * GoAstra CLI - Migration Utilities
 *
 * Helper functions for SQL parsing, version management,
 * driver detection, and database-specific operations.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package migrator

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/***
 * extractSQL parses migration content to extract directional SQL.
 * GoAstra migrations use -- @up and -- @down markers for clarity.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func extractSQL(content, direction string) string {
	upMarker := "-- @up"
	downMarker := "-- @down"

	upIdx := strings.Index(strings.ToLower(content), upMarker)
	downIdx := strings.Index(strings.ToLower(content), downMarker)

	if direction == "up" {
		if upIdx == -1 {
			if downIdx == -1 {
				return content
			}
			return content[:downIdx]
		}
		startIdx := upIdx + len(upMarker)
		if downIdx > upIdx {
			return strings.TrimSpace(content[startIdx:downIdx])
		}
		return strings.TrimSpace(content[startIdx:])
	}

	if downIdx == -1 {
		return ""
	}
	startIdx := downIdx + len(downMarker)
	return strings.TrimSpace(content[startIdx:])
}

/***
 * GenerateVersion creates a GoAstra timestamp-based version identifier.
 * Format: YYYYMMDDHHmmss ensures chronological sortability.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func GenerateVersion() string {
	return time.Now().Format("20060102150405")
}

/***
 * ParseVersion converts a version string to a time.Time value.
 * Enables temporal analysis of migration sequences.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func ParseVersion(version string) (time.Time, error) {
	if len(version) != 14 {
		return time.Time{}, fmt.Errorf("invalid version format: %s", version)
	}

	year, _ := strconv.Atoi(version[0:4])
	month, _ := strconv.Atoi(version[4:6])
	day, _ := strconv.Atoi(version[6:8])
	hour, _ := strconv.Atoi(version[8:10])
	minute, _ := strconv.Atoi(version[10:12])
	second, _ := strconv.Atoi(version[12:14])

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}

/***
 * DetectDriverFromURL attempts to determine the database driver from a connection URL.
 * Supports common URL formats for MySQL, PostgreSQL, and SQLite.
 *
 * Examples:
 *   - mysql://user:pass@localhost/db    -> mysql
 *   - postgres://user:pass@localhost/db -> postgres
 *   - user:pass@tcp(localhost:3306)/db  -> mysql (DSN format)
 *   - file:test.db                      -> sqlite3
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func DetectDriverFromURL(url string) string {
	url = strings.ToLower(url)

	// Check URL scheme prefixes
	if strings.HasPrefix(url, "postgres://") || strings.HasPrefix(url, "postgresql://") {
		return DriverPostgres
	}

	if strings.HasPrefix(url, "mysql://") {
		return DriverMySQL
	}

	if strings.HasPrefix(url, "sqlite://") || strings.HasPrefix(url, "file:") {
		return DriverSQLite
	}

	// Check for MySQL DSN format: user:pass@tcp(host:port)/dbname
	if strings.Contains(url, "@tcp(") || strings.Contains(url, "@unix(") {
		return DriverMySQL
	}

	// Check for SQLite file paths
	if strings.HasSuffix(url, ".db") || strings.HasSuffix(url, ".sqlite") || strings.HasSuffix(url, ".sqlite3") {
		return DriverSQLite
	}

	// Check for PostgreSQL connection string format
	if strings.Contains(url, "host=") && strings.Contains(url, "dbname=") {
		return DriverPostgres
	}

	// Default to MySQL for localhost development
	return DriverMySQL
}

/***
 * dropAllTables removes all tables from the database.
 * Uses database-specific queries optimized for each supported driver.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) dropAllTables() error {
	var query string

	switch m.driver {
	case DriverPostgres:
		query = `
			DO $$ DECLARE
				r RECORD;
			BEGIN
				FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
					EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
				END LOOP;
			END $$;
		`
	case DriverMySQL:
		// MySQL requires multiple statements for dropping all tables
		return m.dropAllTablesMySQL()

	case DriverSQLite:
		query = `
			PRAGMA writable_schema = 1;
			DELETE FROM sqlite_master WHERE type IN ('table', 'index', 'trigger');
			PRAGMA writable_schema = 0;
			VACUUM;
		`
	default:
		return m.dropAllTablesMySQL()
	}

	_, err := m.db.Exec(query)
	return err
}

/***
 * dropAllTablesMySQL handles MySQL-specific table dropping.
 * Disables foreign key checks, drops all tables, then re-enables checks.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) dropAllTablesMySQL() error {
	// Disable foreign key checks
	if _, err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Get all table names
	rows, err := m.db.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
	`)
	if err != nil {
		return fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	// Drop each table
	for _, table := range tables {
		if _, err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Re-enable foreign key checks
	if _, err := m.db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	return nil
}

/***
 * FormatDSN creates a MySQL DSN string from connection parameters.
 * Format: user:password@tcp(host:port)/dbname?parseTime=true
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func FormatDSN(user, password, host, port, dbname string) string {
	if port == "" {
		port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, dbname)
}
