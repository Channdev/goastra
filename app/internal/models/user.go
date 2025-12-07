/*
 * GoAstra Backend - User Model
 *
 * User entity and related data structures.
 */
package models

import "time"

/*
 * User represents an authenticated user in the system.
 */
type User struct {
	ID           uint      `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	Role         string    `db:"role" json:"role"`
	Active       bool      `db:"active" json:"active"`
	Avatar       *string   `db:"avatar" json:"avatar,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

/*
 * UserResponse is the public representation of a User.
 */
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	Avatar    *string   `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
 * ToResponse converts User to its public response representation.
 */
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		Active:    u.Active,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

/*
 * CreateUserInput defines the input for creating a user.
 */
type CreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

/*
 * UpdateUserInput defines the input for updating a user.
 */
type UpdateUserInput struct {
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Role   *string `json:"role,omitempty"`
	Active *bool   `json:"active,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
}

/*
 * LoginInput defines the input for user login.
 */
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

/*
 * AuthResponse is returned after successful authentication.
 */
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
	TokenType    string        `json:"token_type"`
	User         *UserResponse `json:"user"`
}
