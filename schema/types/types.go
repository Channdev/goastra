/*
 * GoAstra Schema Types
 *
 * Define your shared types here. These will be converted to TypeScript
 * using the 'goastra typesync' command.
 *
 * All structs in this package are parsed for TypeScript generation.
 * Use JSON tags to control field names in the generated interfaces.
 */
package types

import "time"

/*
 * BaseModel provides common fields for all database models.
 * Automatically handles ID, timestamps, and soft deletes.
 */
type BaseModel struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
 * User represents an authenticated user in the system.
 */
type User struct {
	BaseModel
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
	Avatar   string `json:"avatar,omitempty"`
}

/*
 * UserProfile represents publicly visible user information.
 */
type UserProfile struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Bio    string `json:"bio,omitempty"`
}

/*
 * PaginatedResponse wraps list responses with pagination metadata.
 */
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

/*
 * APIError represents a standardized error response.
 */
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

/*
 * ValidationError represents a field-level validation error.
 */
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
}

/*
 * AuthTokens contains authentication token pair.
 */
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

/*
 * LoginRequest defines login endpoint payload.
 */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
 * RegisterRequest defines registration endpoint payload.
 */
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

/*
 * HealthStatus represents system health check response.
 */
type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}
