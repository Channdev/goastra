/*
 * GoAstra CLI - Ent Client Template
 *
 * Generates Ent ORM client setup and database connection.
 * Supports MySQL and PostgreSQL with migrations.
 */
package ent

// ClientGoMySQL returns the Ent client template for MySQL.
func ClientGoMySQL() string {
	return `package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"

	"app/ent"
)

/*
 * Client wraps the Ent client with additional functionality.
 */
type Client struct {
	*ent.Client
}

/*
 * Connect establishes a database connection and returns an Ent client.
 */
func Connect() (*Client, error) {
	url := getDatabaseURL()
	if url == "" {
		return nil, nil
	}

	drv, err := sql.Open(dialect.MySQL, url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db := drv.DB()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	client := ent.NewClient(ent.Driver(drv))
	return &Client{Client: client}, nil
}

/*
 * getDatabaseURL builds MySQL connection string from environment.
 */
func getDatabaseURL() string {
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

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
 * Migrate runs auto-migration for the Ent schema.
 */
func (c *Client) Migrate(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Schema.Create(ctx)
}

/*
 * Health checks the database connection.
 */
func (c *Client) Health() error {
	if c == nil || c.Client == nil {
		return nil
	}
	ctx := context.Background()
	// Simple query to check connection
	_, err := c.User.Query().Limit(1).All(ctx)
	if err != nil && err.Error() != "ent: user not found" {
		return err
	}
	return nil
}

/*
 * Close terminates the database connection.
 */
func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Client.Close()
}
`
}

// ClientGoPostgres returns the Ent client template for PostgreSQL.
func ClientGoPostgres() string {
	return `package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"

	"app/ent"
)

/*
 * Client wraps the Ent client with additional functionality.
 */
type Client struct {
	*ent.Client
}

/*
 * Connect establishes a database connection and returns an Ent client.
 */
func Connect() (*Client, error) {
	url := getDatabaseURL()
	if url == "" {
		return nil, nil
	}

	drv, err := sql.Open(dialect.Postgres, url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db := drv.DB()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	client := ent.NewClient(ent.Driver(drv))
	return &Client{Client: client}, nil
}

/*
 * getDatabaseURL builds PostgreSQL connection string from environment.
 */
func getDatabaseURL() string {
	if url := os.Getenv("DB_URL"); url != "" {
		return url
	}

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
 * Migrate runs auto-migration for the Ent schema.
 */
func (c *Client) Migrate(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Schema.Create(ctx)
}

/*
 * Health checks the database connection.
 */
func (c *Client) Health() error {
	if c == nil || c.Client == nil {
		return nil
	}
	ctx := context.Background()
	_, err := c.User.Query().Limit(1).All(ctx)
	if err != nil && err.Error() != "ent: user not found" {
		return err
	}
	return nil
}

/*
 * Close terminates the database connection.
 */
func (c *Client) Close() error {
	if c == nil || c.Client == nil {
		return nil
	}
	return c.Client.Close()
}
`
}

// ClientGo returns the appropriate Ent client template based on driver.
func ClientGo(driver string) string {
	if driver == "mysql" {
		return ClientGoMySQL()
	}
	return ClientGoPostgres()
}
