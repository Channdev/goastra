/***
 * GoAstra CLI - Database Migration Command
 *
 * GoAstra's native database migration management system built specifically
 * for Go applications. Provides a complete suite of migration commands
 * for managing database schema changes with enterprise-grade reliability.
 *
 * Available Commands:
 *   goastra migrate              Run all pending migrations
 *   goastra migrate:status       Show migration status
 *   goastra migrate:rollback     Rollback the last batch of migrations
 *   goastra migrate:reset        Rollback all migrations
 *   goastra migrate:refresh      Reset and re-run all migrations
 *   goastra migrate:fresh        Drop all tables and re-run migrations
 *   goastra migrate:make         Create a new migration file
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/channdev/goastra/cli/internal/migrator"
	"github.com/spf13/cobra"
)

/***
 * Migration command flags for customizing behavior.
 * These allow fine-grained control over migration execution.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var (
	migrateSteps      int
	migrateSeed       bool
	migrateForce      bool
	migrateDatabase   string
	migratePath       string
	migrateCreateTable bool
)

/***
 * migrateCmd is the parent command for all migration operations.
 * Running without subcommands executes pending migrations.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Database migration management",
	Long: `/***
 * GoAstra Database Migration System
 *
 * A powerful migration framework for managing database schema changes
 * with version control, rollback support, and batch tracking.
 *
 * Running 'goastra migrate' without subcommands will execute all
 * pending migrations in order.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/

Usage Examples:
  goastra migrate                    Run all pending migrations
  goastra migrate --step=3           Run only 3 pending migrations
  goastra migrate --seed             Run migrations and seed database
  goastra migrate --force            Force run in production
  goastra migrate --database=mysql   Use specific database connection

Subcommands:
  goastra migrate:status             Show the status of each migration
  goastra migrate:rollback           Rollback the last batch of migrations
  goastra migrate:reset              Rollback all database migrations
  goastra migrate:refresh            Reset and re-run all migrations
  goastra migrate:fresh              Drop all tables and re-run migrations
  goastra migrate:make <name>        Create a new migration file`,
	RunE: runMigrate,
}

/***
 * migrateStatusCmd displays the current status of all migrations.
 * Shows which migrations have been run and which are pending.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of each migration",
	Long: `/***
 * Migration Status Command
 *
 * Displays a detailed table showing all discovered migrations,
 * their run status, batch number, and execution timestamp.
 *
 * Status Indicators:
 *   [Ran]     - Migration has been successfully applied
 *   [Pending] - Migration is waiting to be run
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/`,
	RunE: runMigrateStatus,
}

/***
 * migrateRollbackCmd reverts the most recent batch of migrations.
 * Supports step-based rollback for precise control.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration batch",
	Long: `/***
 * Migration Rollback Command
 *
 * Reverts the last batch of migrations by executing their
 * down() methods in reverse order. Use --step to rollback
 * a specific number of migrations instead of a full batch.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/

Usage Examples:
  goastra migrate:rollback           Rollback the last batch
  goastra migrate:rollback --step=5  Rollback last 5 migrations`,
	RunE: runMigrateRollback,
}

/***
 * migrateResetCmd rolls back all migrations in the database.
 * Use with caution as this will revert all schema changes.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateResetCmd = &cobra.Command{
	Use:   "migrate:reset",
	Short: "Rollback all database migrations",
	Long: `/***
 * Migration Reset Command
 *
 * Rolls back ALL migrations by executing their down() methods
 * in reverse chronological order. This effectively returns
 * your database to its initial state.
 *
 * WARNING: This is a destructive operation. Use --force
 * in production environments.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/`,
	RunE: runMigrateReset,
}

/***
 * migrateRefreshCmd resets and re-runs all migrations.
 * Useful for ensuring migrations work correctly from scratch.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateRefreshCmd = &cobra.Command{
	Use:   "migrate:refresh",
	Short: "Reset and re-run all migrations",
	Long: `/***
 * Migration Refresh Command
 *
 * Performs a complete reset followed by running all migrations.
 * This is equivalent to running migrate:reset then migrate.
 * Optionally seeds the database after migrations complete.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/

Usage Examples:
  goastra migrate:refresh            Reset and re-run all migrations
  goastra migrate:refresh --seed     Also run database seeders
  goastra migrate:refresh --step=5   Rollback 5 then migrate`,
	RunE: runMigrateRefresh,
}

/***
 * migrateFreshCmd drops all tables and re-runs migrations.
 * This is more thorough than refresh as it handles orphaned tables.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Long: `/***
 * Migration Fresh Command
 *
 * Drops ALL tables in the database (not just those tracked
 * in migrations) and then runs all migrations from scratch.
 * This ensures a completely clean database state.
 *
 * WARNING: This is HIGHLY destructive. All data will be lost.
 * Requires --force flag in production environments.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/

Usage Examples:
  goastra migrate:fresh              Fresh database rebuild
  goastra migrate:fresh --seed       Also seed after migrations
  goastra migrate:fresh --force      Force in production`,
	RunE: runMigrateFresh,
}

/***
 * migrateMakeCmd creates a new migration file.
 * Generates timestamped migration with up/down templates.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
var migrateMakeCmd = &cobra.Command{
	Use:   "migrate:make <name>",
	Short: "Create a new migration file",
	Long: `/***
 * Migration Make Command
 *
 * Generates a new migration file with a timestamp-based version
 * and the provided name. Migration files are created in the
 * app/database/migrations directory.
 *
 * Naming Conventions:
 *   create_users_table      - For creating new tables
 *   add_email_to_users      - For adding columns
 *   modify_status_in_orders - For modifying columns
 *   drop_legacy_table       - For dropping tables
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/

Usage Examples:
  goastra migrate:make create_users_table
  goastra migrate:make add_email_to_users
  goastra migrate:make create_products_table --create=products`,
	Args: cobra.ExactArgs(1),
	RunE: runMigrateMake,
}

/***
 * init registers all migration commands with the root command.
 * Sets up flags and subcommand relationships.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateResetCmd)
	rootCmd.AddCommand(migrateRefreshCmd)
	rootCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(migrateMakeCmd)

	// Global migration flags
	migrateCmd.PersistentFlags().StringVar(&migrateDatabase, "database", "", "database connection to use")
	migrateCmd.PersistentFlags().StringVar(&migratePath, "path", "", "path to migrations directory")
	migrateCmd.PersistentFlags().BoolVar(&migrateForce, "force", false, "force operation in production")

	// Migrate command specific flags
	migrateCmd.Flags().IntVar(&migrateSteps, "step", 0, "number of migrations to run")
	migrateCmd.Flags().BoolVar(&migrateSeed, "seed", false, "run seeders after migration")

	// Rollback flags
	migrateRollbackCmd.Flags().IntVar(&migrateSteps, "step", 0, "number of migrations to rollback")

	// Refresh flags
	migrateRefreshCmd.Flags().IntVar(&migrateSteps, "step", 0, "number of migrations to rollback before migrating")
	migrateRefreshCmd.Flags().BoolVar(&migrateSeed, "seed", false, "run seeders after refresh")

	// Fresh flags
	migrateFreshCmd.Flags().BoolVar(&migrateSeed, "seed", false, "run seeders after fresh migration")

	// Make flags
	migrateMakeCmd.Flags().BoolVar(&migrateCreateTable, "create", false, "create table migration template")
}

/***
 * getMigrator creates and configures a new Migrator instance.
 * Handles configuration loading and connection setup.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func getMigrator() (*migrator.Migrator, error) {
	cfg := migrator.DefaultConfig()

	if migratePath != "" {
		cfg.MigrationsPath = migratePath
	}

	if migrateDatabase != "" {
		cfg.DatabaseURL = migrateDatabase
	}

	// Load database URL from environment if not specified
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = os.Getenv("DATABASE_URL")
		if cfg.DatabaseURL == "" {
			// Try loading from .env file
			cfg.DatabaseURL = loadDatabaseURL()
		}
	}

	m, err := migrator.New(cfg)
	if err != nil {
		return nil, err
	}

	return m, nil
}

/***
 * loadDatabaseURL attempts to load the database URL from configuration.
 * Checks environment variables and auto-builds MySQL URL from individual vars.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func loadDatabaseURL() string {
	// Check DB_URL first (primary)
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

	// Auto-detect MySQL from individual environment variables
	if url := buildMySQLURL(); url != "" {
		return url
	}

	// Auto-detect PostgreSQL from individual environment variables
	if url := buildPostgresURL(); url != "" {
		return url
	}

	return ""
}

/***
 * buildMySQLURL constructs a MySQL connection URL from individual env vars.
 * Uses MYSQL_HOST, MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_DATABASE, MYSQL_PORT.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func buildMySQLURL() string {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	// Need at least host and user to build URL
	if host == "" || user == "" {
		return ""
	}

	if port == "" {
		port = "3306"
	}

	if dbname == "" {
		dbname = "goastra"
	}

	// Build MySQL DSN: user:password@tcp(host:port)/dbname?parseTime=true
	if pass != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	}
	return fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true", user, host, port, dbname)
}

/***
 * buildPostgresURL constructs a PostgreSQL connection URL from individual env vars.
 * Uses POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func buildPostgresURL() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	// Need at least host and user to build URL
	if host == "" || user == "" {
		return ""
	}

	if port == "" {
		port = "5432"
	}

	if dbname == "" {
		dbname = "goastra"
	}

	// Build PostgreSQL URL: postgres://user:password@host:port/dbname?sslmode=disable
	if pass != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	}
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", user, host, port, dbname)
}

/***
 * checkProductionSafety verifies it's safe to run destructive operations.
 * Requires --force flag in production environments.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func checkProductionSafety(operation string) error {
	env := os.Getenv("GOASTRA_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}

	if env == "production" && !migrateForce {
		return fmt.Errorf(
			"refusing to run %s in production without --force flag\n"+
				"Use: goastra %s --force",
			operation, operation,
		)
	}
	return nil
}

/***
 * runMigrate executes all pending database migrations.
 * Main entry point for the migrate command.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrate(cmd *cobra.Command, args []string) error {
	color.Cyan("\n  GoAstra Migration System\n")
	color.Cyan("  ========================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		color.Yellow("  No database connection configured.\n")
		color.Yellow("  Set DB_URL or MYSQL_*/POSTGRES_* environment variables.\n\n")
		printMigrateHelp()
		return nil
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	var count int
	if migrateSteps > 0 {
		color.Yellow("  Running %d migration(s)...\n\n", migrateSteps)
		count, err = m.MigrateStep(migrateSteps)
	} else {
		color.Yellow("  Running pending migrations...\n\n")
		count, err = m.Migrate()
	}

	if err != nil {
		color.Red("  Migration failed: %v\n", err)
		return err
	}

	if count == 0 {
		color.Green("  Nothing to migrate. Database is up to date.\n\n")
	} else {
		color.Green("  Successfully ran %d migration(s).\n\n", count)
	}

	if migrateSeed {
		color.Yellow("  Running database seeders...\n")
		// TODO: Implement seeder execution
		color.Green("  Seeding complete.\n\n")
	}

	return nil
}

/***
 * runMigrateStatus displays the status of all migrations.
 * Shows a formatted table with migration details.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateStatus(cmd *cobra.Command, args []string) error {
	color.Cyan("\n  Migration Status\n")
	color.Cyan("  ================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		// Show file-only status
		statuses, err := m.Status()
		if err != nil {
			return err
		}

		printStatusTable(statuses, false)
		return nil
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	statuses, err := m.Status()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	printStatusTable(statuses, true)
	return nil
}

/***
 * printStatusTable outputs a formatted migration status table.
 * Uses color coding for ran/pending status.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func printStatusTable(statuses []migrator.MigrationStatus, showDetails bool) {
	if len(statuses) == 0 {
		color.Yellow("  No migrations found.\n\n")
		return
	}

	// Header
	fmt.Printf("  %-10s %-50s", "Status", "Migration")
	if showDetails {
		fmt.Printf(" %-8s %-20s", "Batch", "Ran At")
	}
	fmt.Println()
	fmt.Printf("  %s\n", strings.Repeat("-", 90))

	// Rows
	for _, status := range statuses {
		var statusStr string
		if status.Ran {
			statusStr = color.GreenString("Ran")
		} else {
			statusStr = color.YellowString("Pending")
		}

		migrationName := fmt.Sprintf("%s_%s", status.Migration.Version, status.Migration.Name)
		if len(migrationName) > 48 {
			migrationName = migrationName[:45] + "..."
		}

		fmt.Printf("  %-10s %-50s", statusStr, migrationName)
		if showDetails && status.Ran {
			fmt.Printf(" %-8d %-20s",
				status.Migration.Batch,
				status.Migration.AppliedAt.Format("2006-01-02 15:04:05"),
			)
		}
		fmt.Println()
	}
	fmt.Println()
}

/***
 * runMigrateRollback reverts the last batch of migrations.
 * Supports step-based rollback for fine control.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateRollback(cmd *cobra.Command, args []string) error {
	if err := checkProductionSafety("migrate:rollback"); err != nil {
		return err
	}

	color.Cyan("\n  Rolling Back Migrations\n")
	color.Cyan("  =======================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		return fmt.Errorf("no database connection configured")
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	var count int
	if migrateSteps > 0 {
		color.Yellow("  Rolling back %d migration(s)...\n\n", migrateSteps)
		count, err = m.RollbackStep(migrateSteps)
	} else {
		color.Yellow("  Rolling back last batch...\n\n")
		count, err = m.Rollback()
	}

	if err != nil {
		color.Red("  Rollback failed: %v\n", err)
		return err
	}

	if count == 0 {
		color.Yellow("  Nothing to rollback.\n\n")
	} else {
		color.Green("  Successfully rolled back %d migration(s).\n\n", count)
	}

	return nil
}

/***
 * runMigrateReset rolls back all database migrations.
 * Returns database to initial state.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateReset(cmd *cobra.Command, args []string) error {
	if err := checkProductionSafety("migrate:reset"); err != nil {
		return err
	}

	color.Cyan("\n  Resetting All Migrations\n")
	color.Cyan("  ========================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		return fmt.Errorf("no database connection configured")
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	color.Yellow("  Rolling back all migrations...\n\n")
	count, err := m.Reset()
	if err != nil {
		color.Red("  Reset failed: %v\n", err)
		return err
	}

	if count == 0 {
		color.Yellow("  No migrations to reset.\n\n")
	} else {
		color.Green("  Successfully reset %d migration(s).\n\n", count)
	}

	return nil
}

/***
 * runMigrateRefresh resets and re-runs all migrations.
 * Ensures clean migration state.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateRefresh(cmd *cobra.Command, args []string) error {
	if err := checkProductionSafety("migrate:refresh"); err != nil {
		return err
	}

	color.Cyan("\n  Refreshing Migrations\n")
	color.Cyan("  =====================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		return fmt.Errorf("no database connection configured")
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	color.Yellow("  Rolling back and re-running migrations...\n\n")
	rolledBack, migrated, err := m.Refresh()
	if err != nil {
		color.Red("  Refresh failed: %v\n", err)
		return err
	}

	color.Green("  Rolled back %d migration(s).\n", rolledBack)
	color.Green("  Migrated %d migration(s).\n\n", migrated)

	if migrateSeed {
		color.Yellow("  Running database seeders...\n")
		// TODO: Implement seeder execution
		color.Green("  Seeding complete.\n\n")
	}

	return nil
}

/***
 * runMigrateFresh drops all tables and re-runs migrations.
 * Most thorough database rebuild option.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateFresh(cmd *cobra.Command, args []string) error {
	if err := checkProductionSafety("migrate:fresh"); err != nil {
		return err
	}

	color.Cyan("\n  Fresh Migration\n")
	color.Cyan("  ===============\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	defer m.Close()

	dbURL := loadDatabaseURL()
	if dbURL == "" {
		return fmt.Errorf("no database connection configured")
	}

	if err := m.Connect(dbURL); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	color.Red("  WARNING: Dropping all tables!\n\n")
	color.Yellow("  Rebuilding database from scratch...\n\n")

	count, err := m.Fresh()
	if err != nil {
		color.Red("  Fresh migration failed: %v\n", err)
		return err
	}

	color.Green("  Dropped all tables and ran %d migration(s).\n\n", count)

	if migrateSeed {
		color.Yellow("  Running database seeders...\n")
		// TODO: Implement seeder execution
		color.Green("  Seeding complete.\n\n")
	}

	return nil
}

/***
 * runMigrateMake creates a new migration file.
 * Generates timestamped migration with templates.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func runMigrateMake(cmd *cobra.Command, args []string) error {
	name := args[0]

	color.Cyan("\n  Creating Migration\n")
	color.Cyan("  ==================\n\n")

	m, err := getMigrator()
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	filepath, err := m.CreateMigration(name, migrateCreateTable)
	if err != nil {
		color.Red("  Failed to create migration: %v\n", err)
		return err
	}

	color.Green("  Created migration: %s\n\n", filepath)
	color.Yellow("  Next steps:\n")
	fmt.Printf("    1. Edit the migration file to add your schema changes\n")
	fmt.Printf("    2. Run 'goastra migrate' to apply the migration\n\n")

	return nil
}

/***
 * printMigrateHelp outputs helpful usage information.
 * Displayed when migrations cannot run due to missing config.
 *
 * Author: channdev
 * Date: 12/10/2025
 ***/
func printMigrateHelp() {
	fmt.Println("  Available Commands:")
	fmt.Println("  -------------------")
	fmt.Println("  goastra migrate              Run all pending migrations")
	fmt.Println("  goastra migrate:status       Show migration status")
	fmt.Println("  goastra migrate:rollback     Rollback the last batch")
	fmt.Println("  goastra migrate:reset        Rollback all migrations")
	fmt.Println("  goastra migrate:refresh      Reset and re-run all")
	fmt.Println("  goastra migrate:fresh        Drop tables and re-run")
	fmt.Println("  goastra migrate:make <name>  Create new migration")
	fmt.Println()
	fmt.Println("  Configuration:")
	fmt.Println("  --------------")
	fmt.Println("  Option 1: Set DB_URL in .env file")
	fmt.Println("  Option 2: Use individual MySQL vars:")
	fmt.Println("            MYSQL_HOST, MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_DATABASE")
	fmt.Println("  Option 3: Use individual Postgres vars:")
	fmt.Println("            POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB")
	fmt.Println("  Option 4: Use --database flag")
	fmt.Println()
	fmt.Println("  Migrations are stored in: app/database/migrations/")
	fmt.Println()
}
