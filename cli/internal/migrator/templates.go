/***
 * GoAstra CLI - Migration Templates
 *
 * SQL template generation for creating new migration files.
 * Provides standardized templates for table creation and
 * custom schema modifications with multi-driver support.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/***
 * CreateMigration generates a new migration file with GoAstra conventions.
 * Uses timestamp-based versioning (YYYYMMDDHHmmss) for guaranteed ordering.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) CreateMigration(name string, createTable bool) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	safeName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	filename := fmt.Sprintf("%s_%s.sql", timestamp, safeName)
	migrationPath := filepath.Join(m.migrationsPath, filename)

	if err := os.MkdirAll(m.migrationsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	var content string
	if createTable {
		tableName := extractTableName(safeName)
		content = m.generateCreateTableMigration(tableName)
	} else {
		content = generateBlankMigration(safeName)
	}

	if err := os.WriteFile(migrationPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write migration file: %w", err)
	}

	return migrationPath, nil
}

/***
 * generateCreateTableMigration produces a driver-specific table creation template.
 * Includes standard columns: id, created_at, and updated_at.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) generateCreateTableMigration(tableName string) string {
	switch m.driver {
	case DriverMySQL:
		return generateMySQLTableMigration(tableName)
	case DriverPostgres:
		return generatePostgresTableMigration(tableName)
	case DriverSQLite:
		return generateSQLiteTableMigration(tableName)
	default:
		return generateMySQLTableMigration(tableName)
	}
}

/***
 * generateMySQLTableMigration creates a MySQL-compatible table template.
 * Uses INT AUTO_INCREMENT and DATETIME types.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func generateMySQLTableMigration(tableName string) string {
	return fmt.Sprintf(`-- GoAstra Migration
-- Table: %s
-- Driver: MySQL

-- @up
CREATE TABLE %s (
    id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- @down
DROP TABLE IF EXISTS %s;
`, tableName, tableName, tableName)
}

/***
 * generatePostgresTableMigration creates a PostgreSQL-compatible table template.
 * Uses SERIAL and TIMESTAMP WITH TIME ZONE types.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func generatePostgresTableMigration(tableName string) string {
	return fmt.Sprintf(`-- GoAstra Migration
-- Table: %s
-- Driver: PostgreSQL

-- @up
CREATE TABLE %s (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- @down
DROP TABLE IF EXISTS %s;
`, tableName, tableName, tableName)
}

/***
 * generateSQLiteTableMigration creates a SQLite-compatible table template.
 * Uses INTEGER PRIMARY KEY AUTOINCREMENT and DATETIME types.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func generateSQLiteTableMigration(tableName string) string {
	return fmt.Sprintf(`-- GoAstra Migration
-- Table: %s
-- Driver: SQLite

-- @up
CREATE TABLE %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- @down
DROP TABLE IF EXISTS %s;
`, tableName, tableName, tableName)
}

/***
 * generateBlankMigration creates an empty GoAstra migration template.
 * Provides the standard structure for custom schema modifications.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func generateBlankMigration(name string) string {
	return fmt.Sprintf(`-- GoAstra Migration: %s
-- Created: %s

-- @up
-- Add your forward migration SQL here

-- @down
-- Add your rollback migration SQL here
`, name, time.Now().Format("2006-01-02 15:04:05"))
}

/***
 * extractTableName derives the table name from a migration identifier.
 * Handles GoAstra naming conventions (e.g., create_users_table -> users).
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func extractTableName(name string) string {
	name = strings.TrimPrefix(name, "create_")
	name = strings.TrimSuffix(name, "_table")
	return name
}
