/*
 * GoAstra CLI - SQLx Repository Template
 *
 * Generates repository pattern templates for SQLx.
 * Provides base CRUD operations with pagination support.
 */
package sqlx

// RepositoryGo returns the repository.go template.
func RepositoryGo() string {
	return `package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"app/internal/database"
)

/*
 * Repository is a generic interface for data access operations.
 */
type Repository[T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

/*
 * BaseRepository provides common repository utilities using SQLx.
 */
type BaseRepository struct {
	db *database.DB
}

/*
 * NewBaseRepository creates a new base repository.
 */
func NewBaseRepository(db *database.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

/*
 * DB returns the underlying database connection.
 */
func (r *BaseRepository) DB() *sqlx.DB {
	if r.db == nil {
		return nil
	}
	return r.db.DB
}

/*
 * Transaction executes a function within a database transaction.
 * Automatically commits on success or rolls back on error.
 */
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/*
 * QueryRow executes a query that returns a single row.
 */
func (r *BaseRepository) QueryRow(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.db.GetContext(ctx, dest, query, args...)
}

/*
 * Query executes a query that returns multiple rows.
 */
func (r *BaseRepository) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return r.db.SelectContext(ctx, dest, query, args...)
}

/*
 * Exec executes a query that doesn't return rows.
 */
func (r *BaseRepository) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

/*
 * NamedExec executes a named query that doesn't return rows.
 */
func (r *BaseRepository) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return r.db.NamedExecContext(ctx, query, arg)
}

/*
 * PaginationParams holds pagination query parameters.
 */
type PaginationParams struct {
	Page     int
	PageSize int
}

/*
 * DefaultPagination returns default pagination settings.
 */
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 20,
	}
}

/*
 * Offset calculates the database offset for pagination.
 */
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

/*
 * Limit returns the page size for limit clause.
 */
func (p PaginationParams) Limit() int {
	return p.PageSize
}
`
}
