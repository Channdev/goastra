/*
 * GoAstra CLI - REST Handlers Template
 *
 * Generates HTTP handler templates for REST APIs.
 * Provides base handler patterns and health check endpoints.
 */
package rest

// HandlersGo returns the handlers.go template.
func HandlersGo() string {
	return `package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * HealthHandler handles health check endpoints.
 */
type HealthHandler struct{}

/*
 * NewHealthHandler creates a new health handler instance.
 */
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

/*
 * Health returns the service health status.
 */
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Service is running",
		"version": "1.0.0",
	})
}

/*
 * Ready returns the service readiness status.
 */
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

/*
 * BaseHandler provides common handler utilities.
 */
type BaseHandler struct{}

/*
 * OK sends a success response with data.
 */
func (h *BaseHandler) OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

/*
 * Created sends a 201 response with the created resource.
 */
func (h *BaseHandler) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

/*
 * NoContent sends a 204 response.
 */
func (h *BaseHandler) NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

/*
 * BadRequest sends a 400 error response.
 */
func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": message,
		"code":  "BAD_REQUEST",
	})
}

/*
 * Unauthorized sends a 401 error response.
 */
func (h *BaseHandler) Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": message,
		"code":  "UNAUTHORIZED",
	})
}

/*
 * Forbidden sends a 403 error response.
 */
func (h *BaseHandler) Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{
		"error": message,
		"code":  "FORBIDDEN",
	})
}

/*
 * NotFound sends a 404 error response.
 */
func (h *BaseHandler) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": message,
		"code":  "NOT_FOUND",
	})
}

/*
 * InternalError sends a 500 error response.
 */
func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
		"code":  "INTERNAL_ERROR",
	})
}

/*
 * ValidationError sends a 422 error response with validation details.
 */
func (h *BaseHandler) ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"error":   "Validation failed",
		"code":    "VALIDATION_ERROR",
		"details": errors,
	})
}
`
}

// ServicesGo returns the services.go template.
func ServicesGo() string {
	return `package services

import "context"

/*
 * Service is a generic interface for business logic operations.
 */
type Service[T any] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
}

/*
 * BaseService provides common service utilities.
 */
type BaseService struct{}

/*
 * NewBaseService creates a new base service instance.
 */
func NewBaseService() *BaseService {
	return &BaseService{}
}
`
}

// RepositoryGo returns the repository.go template.
func RepositoryGo() string {
	return `package repository

import "context"

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
 * BaseRepository provides common repository utilities.
 */
type BaseRepository struct{}

/*
 * NewBaseRepository creates a new base repository instance.
 */
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{}
}
`
}
