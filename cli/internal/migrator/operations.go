/***
 * GoAstra CLI - Migration Operations
 *
 * Contains all migration execution operations including migrate,
 * rollback, reset, refresh, and fresh commands with multi-driver support.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package migrator

import (
	"fmt"
)

/***
 * Migrate executes all pending migrations in version order.
 * Returns the count of successfully applied migrations.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Migrate() (int, error) {
	pending, err := m.GetPendingMigrations()
	if err != nil {
		return 0, err
	}

	if len(pending) == 0 {
		return 0, nil
	}

	batch, err := m.GetNextBatch()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, mig := range pending {
		if err := m.runMigration(mig, batch, "up"); err != nil {
			return count, fmt.Errorf("migration %s failed: %w", mig.Name, err)
		}
		count++
	}

	return count, nil
}

/***
 * MigrateStep executes a specified number of pending migrations.
 * Provides controlled, incremental schema updates.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) MigrateStep(steps int) (int, error) {
	pending, err := m.GetPendingMigrations()
	if err != nil {
		return 0, err
	}

	if len(pending) == 0 {
		return 0, nil
	}

	if steps > len(pending) {
		steps = len(pending)
	}

	batch, err := m.GetNextBatch()
	if err != nil {
		return 0, err
	}

	count := 0
	for i := 0; i < steps; i++ {
		if err := m.runMigration(pending[i], batch, "up"); err != nil {
			return count, fmt.Errorf("migration %s failed: %w", pending[i].Name, err)
		}
		count++
	}

	return count, nil
}

/***
 * Rollback reverts the most recent batch of migrations.
 * Executes the down migration for each migration in the batch.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Rollback() (int, error) {
	batch, err := m.GetLastBatch()
	if err != nil {
		return 0, err
	}

	if batch == 0 {
		return 0, nil
	}

	return m.RollbackBatch(batch)
}

/***
 * RollbackBatch reverts all migrations within a specific batch.
 * Migrations are processed in reverse version order.
 * Uses driver-specific placeholder syntax.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) RollbackBatch(batch int) (int, error) {
	query := m.buildBatchSelectQuery()

	rows, err := m.db.Query(query, batch)
	if err != nil {
		return 0, fmt.Errorf("failed to query batch migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var mig Migration
		if err := rows.Scan(&mig.Version, &mig.Name); err != nil {
			return 0, fmt.Errorf("failed to scan migration: %w", err)
		}
		mig.Filename = m.getMigrationFilename(mig.Version, mig.Name)
		migrations = append(migrations, mig)
	}

	count := 0
	for _, mig := range migrations {
		if err := m.runMigration(mig, batch, "down"); err != nil {
			return count, fmt.Errorf("rollback of %s failed: %w", mig.Name, err)
		}
		count++
	}

	return count, nil
}

/***
 * buildBatchSelectQuery creates a driver-specific query for batch selection.
 * Handles placeholder syntax differences between drivers.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) buildBatchSelectQuery() string {
	switch m.driver {
	case DriverPostgres:
		return fmt.Sprintf(`
			SELECT version, name
			FROM %s
			WHERE batch = $1
			ORDER BY version DESC
		`, m.tableName)
	default:
		// MySQL, SQLite use ? placeholders
		return fmt.Sprintf(`
			SELECT version, name
			FROM %s
			WHERE batch = ?
			ORDER BY version DESC
		`, m.tableName)
	}
}

/***
 * RollbackStep reverts a specified number of migrations.
 * Enables precise, controlled schema rollback operations.
 * Uses driver-specific LIMIT syntax.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) RollbackStep(steps int) (int, error) {
	query := m.buildStepSelectQuery()

	rows, err := m.db.Query(query, steps)
	if err != nil {
		return 0, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var mig Migration
		if err := rows.Scan(&mig.Version, &mig.Name, &mig.Batch); err != nil {
			return 0, fmt.Errorf("failed to scan migration: %w", err)
		}
		mig.Filename = m.getMigrationFilename(mig.Version, mig.Name)
		migrations = append(migrations, mig)
	}

	count := 0
	for _, mig := range migrations {
		if err := m.runMigration(mig, mig.Batch, "down"); err != nil {
			return count, fmt.Errorf("rollback of %s failed: %w", mig.Name, err)
		}
		count++
	}

	return count, nil
}

/***
 * buildStepSelectQuery creates a driver-specific query for step-based selection.
 * Handles placeholder and LIMIT syntax differences.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) buildStepSelectQuery() string {
	switch m.driver {
	case DriverPostgres:
		return fmt.Sprintf(`
			SELECT version, name, batch
			FROM %s
			ORDER BY version DESC
			LIMIT $1
		`, m.tableName)
	default:
		// MySQL, SQLite use ? placeholders
		return fmt.Sprintf(`
			SELECT version, name, batch
			FROM %s
			ORDER BY version DESC
			LIMIT ?
		`, m.tableName)
	}
}

/***
 * Reset reverts all applied migrations to initial state.
 * Returns the total count of rolled back migrations.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Reset() (int, error) {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return 0, err
	}

	// Reverse order for rollback
	for i, j := 0, len(applied)-1; i < j; i, j = i+1, j-1 {
		applied[i], applied[j] = applied[j], applied[i]
	}

	count := 0
	for _, mig := range applied {
		mig.Filename = m.getMigrationFilename(mig.Version, mig.Name)
		if err := m.runMigration(mig, mig.Batch, "down"); err != nil {
			return count, fmt.Errorf("rollback of %s failed: %w", mig.Name, err)
		}
		count++
	}

	return count, nil
}

/***
 * Fresh performs a complete database reconstruction.
 * Drops all tables and re-executes all migrations from scratch.
 * WARNING: This is destructive - use only in development environments.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Fresh() (int, error) {
	if err := m.dropAllTables(); err != nil {
		return 0, fmt.Errorf("failed to drop tables: %w", err)
	}

	if err := m.EnsureMigrationTable(); err != nil {
		return 0, err
	}

	return m.Migrate()
}

/***
 * Refresh performs a complete migration cycle reset.
 * Rolls back all migrations then re-applies them in order.
 * Useful for validating migration integrity during development.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Refresh() (int, int, error) {
	rolledBack, err := m.Reset()
	if err != nil {
		return rolledBack, 0, err
	}

	migrated, err := m.Migrate()
	return rolledBack, migrated, err
}

/***
 * Status provides a comprehensive view of all migration states.
 * Returns detailed information about executed and pending migrations.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func (m *Migrator) Status() ([]MigrationStatus, error) {
	allMigrations, err := m.DiscoverMigrations()
	if err != nil {
		return nil, err
	}

	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]Migration)
	for _, mig := range applied {
		appliedMap[mig.Version] = mig
	}

	var statuses []MigrationStatus
	for _, mig := range allMigrations {
		status := MigrationStatus{
			Migration: mig,
			Pending:   true,
			Ran:       false,
		}

		if appliedMig, exists := appliedMap[mig.Version]; exists {
			status.Migration.Batch = appliedMig.Batch
			status.Migration.AppliedAt = appliedMig.AppliedAt
			status.Pending = false
			status.Ran = true
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
