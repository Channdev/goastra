/*
 * GoAstra CLI - SQLx Database Template
 *
 * Generates database connection templates using SQLx.
 * Supports MySQL and PostgreSQL with auto-detection from env vars.
 */
package sqlx

// DatabaseGoMySQL returns the MySQL database.go template.
func DatabaseGoMySQL() string {
	return `package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

/*
 * DB wraps sqlx.DB with additional functionality.
 */
type DB struct {
	*sqlx.DB
}

/*
 * Connect establishes a database connection using environment variables.
 * Returns nil if no database is configured.
 */
func Connect() (*DB, error) {
	url := getDatabaseURL()
	if url == "" {
		return nil, nil
	}

	db, err := sqlx.Connect("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}

/*
 * getDatabaseURL builds MySQL connection string from environment.
 */
func getDatabaseURL() string {
	// Check DB_URL first
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

	// Build from individual vars
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	if host == "" || user == "" {
		return ""
	}
	if port == "" {
		port = "3306"
	}
	if dbname == "" {
		dbname = "goastra"
	}

	if pass != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, pass, host, port, dbname)
	}
	return fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, host, port, dbname)
}

/*
 * Health checks the database connection.
 */
func (db *DB) Health() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.Ping()
}

/*
 * Close terminates the database connection.
 */
func (db *DB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}
`
}

// DatabaseGoPostgres returns the PostgreSQL database.go template.
func DatabaseGoPostgres() string {
	return `package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
 * DB wraps sqlx.DB with additional functionality.
 */
type DB struct {
	*sqlx.DB
}

/*
 * Connect establishes a database connection using environment variables.
 * Returns nil if no database is configured.
 */
func Connect() (*DB, error) {
	url := getDatabaseURL()
	if url == "" {
		return nil, nil
	}

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}

/*
 * getDatabaseURL builds PostgreSQL connection string from environment.
 */
func getDatabaseURL() string {
	// Check DB_URL first
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

	// Build from individual vars
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" {
		return ""
	}
	if port == "" {
		port = "5432"
	}
	if dbname == "" {
		dbname = "goastra"
	}

	if pass != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbname)
	}
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", user, host, port, dbname)
}

/*
 * Health checks the database connection.
 */
func (db *DB) Health() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.Ping()
}

/*
 * Close terminates the database connection.
 */
func (db *DB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.DB.Close()
}
`
}

// DatabaseGo returns the appropriate database template based on driver.
func DatabaseGo(driver string) string {
	if driver == "mysql" {
		return DatabaseGoMySQL()
	}
	return DatabaseGoPostgres()
}
