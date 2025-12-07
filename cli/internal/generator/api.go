/*
 * GoAstra CLI - API Generator
 *
 * Generates REST API components including handlers, services,
 * repositories, and route registrations.
 */
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

/*
 * APIGenerator handles backend API code generation.
 */
type APIGenerator struct {
	name       string
	pascalName string
	camelName  string
}

/*
 * NewAPIGenerator creates a new API generator instance.
 */
func NewAPIGenerator(name string) *APIGenerator {
	return &APIGenerator{
		name:       name,
		pascalName: toPascalCase(name),
		camelName:  toCamelCase(name),
	}
}

/*
 * GenerateHandler creates the HTTP handler file.
 */
func (g *APIGenerator) GenerateHandler() error {
	content := fmt.Sprintf(`/*
 * %s Handler
 *
 * HTTP handlers for %s CRUD operations.
 */
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
 * %sHandler manages HTTP requests for %s resources.
 */
type %sHandler struct {
	service *%sService
}

/*
 * New%sHandler creates a new handler instance.
 */
func New%sHandler(service *%sService) *%sHandler {
	return &%sHandler{service: service}
}

/*
 * List returns all %s resources with pagination.
 */
func (h *%sHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

/*
 * Get returns a single %s by ID.
 */
func (h *%sHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

/*
 * Create adds a new %s resource.
 */
func (h *%sHandler) Create(c *gin.Context) {
	var input Create%sInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Create(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

/*
 * Update modifies an existing %s resource.
 */
func (h *%sHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input Update%sInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Update(c.Request.Context(), uint(id), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

/*
 * Delete removes a %s resource.
 */
func (h *%sHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

/*
 * Create%sInput defines the request body for creating a %s.
 */
type Create%sInput struct {
	/* TODO: Define create input fields */
}

/*
 * Update%sInput defines the request body for updating a %s.
 */
type Update%sInput struct {
	/* TODO: Define update input fields */
}
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName,
		g.name,
		g.pascalName,
		g.name,
		g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
	)

	path := filepath.Join("app/internal/handlers", g.name+"_handler.go")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateService creates the service layer file.
 */
func (g *APIGenerator) GenerateService() error {
	content := fmt.Sprintf(`/*
 * %s Service
 *
 * Business logic for %s operations.
 */
package services

import (
	"context"
)

/*
 * %sService handles business logic for %s resources.
 */
type %sService struct {
	repo *%sRepository
}

/*
 * New%sService creates a new service instance.
 */
func New%sService(repo *%sRepository) *%sService {
	return &%sService{repo: repo}
}

/*
 * List returns paginated %s resources.
 */
func (s *%sService) List(ctx context.Context, page, pageSize int) (*PaginatedResult, error) {
	offset := (page - 1) * pageSize

	items, err := s.repo.FindAll(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &PaginatedResult{
		Data:       items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}, nil
}

/*
 * GetByID retrieves a single %s by ID.
 */
func (s *%sService) GetByID(ctx context.Context, id uint) (*%s, error) {
	return s.repo.FindByID(ctx, id)
}

/*
 * Create adds a new %s resource.
 */
func (s *%sService) Create(ctx context.Context, input *Create%sInput) (*%s, error) {
	entity := &%s{
		/* TODO: Map input fields */
	}

	return s.repo.Create(ctx, entity)
}

/*
 * Update modifies an existing %s resource.
 */
func (s *%sService) Update(ctx context.Context, id uint, input *Update%sInput) (*%s, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	/* TODO: Map input fields to entity */

	return s.repo.Update(ctx, entity)
}

/*
 * Delete removes a %s resource.
 */
func (s *%sService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

/*
 * PaginatedResult wraps list results with pagination metadata.
 */
type PaginatedResult struct {
	Data       interface{}
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName,
		g.name,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName, g.pascalName, g.pascalName,
		g.name,
		g.pascalName,
	)

	path := filepath.Join("app/internal/services", g.name+"_service.go")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateRepository creates the data access layer file.
 */
func (g *APIGenerator) GenerateRepository() error {
	content := fmt.Sprintf(`/*
 * %s Repository
 *
 * Data access layer for %s entities.
 */
package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

/*
 * %sRepository handles database operations for %s.
 */
type %sRepository struct {
	db *sqlx.DB
}

/*
 * New%sRepository creates a new repository instance.
 */
func New%sRepository(db *sqlx.DB) *%sRepository {
	return &%sRepository{db: db}
}

/*
 * FindAll retrieves paginated %s records.
 */
func (r *%sRepository) FindAll(ctx context.Context, offset, limit int) ([]*%s, error) {
	var items []*%s

	query := "SELECT * FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	if err := r.db.SelectContext(ctx, &items, query, limit, offset); err != nil {
		return nil, err
	}

	return items, nil
}

/*
 * FindByID retrieves a single %s by ID.
 */
func (r *%sRepository) FindByID(ctx context.Context, id uint) (*%s, error) {
	var item %s

	query := "SELECT * FROM %s WHERE id = $1"

	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		return nil, err
	}

	return &item, nil
}

/*
 * Create inserts a new %s record.
 */
func (r *%sRepository) Create(ctx context.Context, entity *%s) (*%s, error) {
	query := ` + "`" + `
		INSERT INTO %s (created_at, updated_at)
		VALUES (NOW(), NOW())
		RETURNING *
	` + "`" + `

	if err := r.db.GetContext(ctx, entity, query); err != nil {
		return nil, err
	}

	return entity, nil
}

/*
 * Update modifies an existing %s record.
 */
func (r *%sRepository) Update(ctx context.Context, entity *%s) (*%s, error) {
	query := ` + "`" + `
		UPDATE %s
		SET updated_at = NOW()
		WHERE id = $1
		RETURNING *
	` + "`" + `

	if err := r.db.GetContext(ctx, entity, query, entity.ID); err != nil {
		return nil, err
	}

	return entity, nil
}

/*
 * Delete removes a %s record.
 */
func (r *%sRepository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM %s WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

/*
 * Count returns the total number of %s records.
 */
func (r *%sRepository) Count(ctx context.Context) (int, error) {
	var count int

	query := "SELECT COUNT(*) FROM %s"

	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}

	return count, nil
}
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.pascalName,
		g.pascalName,
		g.name,
		g.pascalName, g.pascalName,
		g.pascalName,
		toSnakeCase(g.name),
		g.name,
		g.pascalName, g.pascalName,
		g.pascalName,
		toSnakeCase(g.name),
		g.name,
		g.pascalName, g.pascalName, g.pascalName,
		toSnakeCase(g.name),
		g.name,
		g.pascalName, g.pascalName, g.pascalName,
		toSnakeCase(g.name),
		g.name,
		g.pascalName,
		toSnakeCase(g.name),
		g.name,
		g.pascalName,
		toSnakeCase(g.name),
	)

	path := filepath.Join("app/internal/repository", g.name+"_repository.go")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateRoutes creates route registration for the API.
 */
func (g *APIGenerator) GenerateRoutes() error {
	content := fmt.Sprintf(`/*
 * %s Routes
 *
 * Route registration for %s API endpoints.
 */
package router

import (
	"github.com/gin-gonic/gin"
)

/*
 * Register%sRoutes adds %s routes to the router.
 */
func Register%sRoutes(router *gin.RouterGroup, handler *%sHandler) {
	%s := router.Group("/%s")
	{
		%s.GET("", handler.List)
		%s.GET("/:id", handler.Get)
		%s.POST("", handler.Create)
		%s.PUT("/:id", handler.Update)
		%s.DELETE("/:id", handler.Delete)
	}
}
`,
		g.pascalName, g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.camelName, toPlural(g.name),
		g.camelName,
		g.camelName,
		g.camelName,
		g.camelName,
		g.camelName,
	)

	path := filepath.Join("app/internal/router", g.name+"_routes.go")
	return os.WriteFile(path, []byte(content), 0644)
}

/* Helper functions for name transformations */

func toPascalCase(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "-", " "), " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(string(pascal[0])) + pascal[1:]
	}
	return pascal
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == '-' {
			result.WriteRune('_')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toPlural(s string) string {
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}
