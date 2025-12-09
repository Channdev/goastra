/*
 * GoAstra CLI - Models Template
 *
 * Generates base data models and DTOs for the application.
 * Includes User model and authentication request/response types.
 */
package common

// ModelsGo returns the models.go template with base entities and DTOs.
func ModelsGo() string {
	return `package models

import "time"

/*
 * User represents a user account in the system.
 */
type User struct {
	ID        uint      ` + "`json:\"id\" db:\"id\"`" + `
	Email     string    ` + "`json:\"email\" db:\"email\"`" + `
	Password  string    ` + "`json:\"-\" db:\"password\"`" + `
	Name      string    ` + "`json:\"name\" db:\"name\"`" + `
	Role      string    ` + "`json:\"role\" db:\"role\"`" + `
	Active    bool      ` + "`json:\"active\" db:\"active\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
}

/*
 * LoginRequest represents the login endpoint request body.
 */
type LoginRequest struct {
	Email    string ` + "`json:\"email\" binding:\"required,email\"`" + `
	Password string ` + "`json:\"password\" binding:\"required,min=6\"`" + `
}

/*
 * RegisterRequest represents the registration endpoint request body.
 */
type RegisterRequest struct {
	Email    string ` + "`json:\"email\" binding:\"required,email\"`" + `
	Password string ` + "`json:\"password\" binding:\"required,min=6\"`" + `
	Name     string ` + "`json:\"name\" binding:\"required,min=2\"`" + `
}

/*
 * AuthResponse represents the authentication response.
 */
type AuthResponse struct {
	Token        string ` + "`json:\"token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token,omitempty\"`" + `
	ExpiresAt    int64  ` + "`json:\"expires_at\"`" + `
	User         *User  ` + "`json:\"user\"`" + `
}

/*
 * RefreshRequest represents the token refresh request body.
 */
type RefreshRequest struct {
	RefreshToken string ` + "`json:\"refresh_token\" binding:\"required\"`" + `
}

/*
 * UpdateUserRequest represents the user update request body.
 */
type UpdateUserRequest struct {
	Email string ` + "`json:\"email\" binding:\"omitempty,email\"`" + `
	Name  string ` + "`json:\"name\" binding:\"omitempty,min=2\"`" + `
	Role  string ` + "`json:\"role\" binding:\"omitempty\"`" + `
}

/*
 * PaginationParams holds pagination query parameters.
 */
type PaginationParams struct {
	Page     int ` + "`form:\"page\" binding:\"omitempty,min=1\"`" + `
	PageSize int ` + "`form:\"page_size\" binding:\"omitempty,min=1,max=100\"`" + `
}

/*
 * PaginatedResponse wraps paginated data with metadata.
 */
type PaginatedResponse struct {
	Data       interface{} ` + "`json:\"data\"`" + `
	Total      int64       ` + "`json:\"total\"`" + `
	Page       int         ` + "`json:\"page\"`" + `
	PageSize   int         ` + "`json:\"page_size\"`" + `
	TotalPages int         ` + "`json:\"total_pages\"`" + `
}

/*
 * NewPaginatedResponse creates a paginated response with calculated metadata.
 */
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
`
}

// ValidatorGo returns the validator.go template for custom validations.
func ValidatorGo() string {
	return `package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

/*
 * Setup initializes custom validators for Gin.
 * Call this during application startup.
 */
func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidations(v)
	}
}

func registerCustomValidations(v *validator.Validate) {
	// Username validation: alphanumeric, underscores, 3-20 chars
	v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		if len(username) < 3 || len(username) > 20 {
			return false
		}
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
		return matched
	})

	// Slug validation: lowercase, hyphens, alphanumeric
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		slug := fl.Field().String()
		matched, _ := regexp.MatchString("^[a-z0-9]+(?:-[a-z0-9]+)*$", slug)
		return matched
	})
}
`
}
