/*
 * GoAstra Backend - Password Hashing
 *
 * Secure password hashing using bcrypt.
 * Provides utilities for password validation and comparison.
 */
package auth

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	/* bcrypt cost factor - higher is more secure but slower */
	DefaultCost = 12
	MinCost     = 10
	MaxCost     = 14
)

/*
 * PasswordHasher handles password hashing operations.
 */
type PasswordHasher struct {
	cost int
}

/*
 * NewPasswordHasher creates a new hasher with default cost.
 */
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{cost: DefaultCost}
}

/*
 * NewPasswordHasherWithCost creates a hasher with specific cost.
 */
func NewPasswordHasherWithCost(cost int) *PasswordHasher {
	if cost < MinCost {
		cost = MinCost
	}
	if cost > MaxCost {
		cost = MaxCost
	}
	return &PasswordHasher{cost: cost}
}

/*
 * Hash generates a bcrypt hash of the password.
 */
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

/*
 * Compare verifies a password against its hash.
 */
func (h *PasswordHasher) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
 * PasswordPolicy defines password requirements.
 */
type PasswordPolicy struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

/*
 * DefaultPasswordPolicy returns standard security requirements.
 */
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		MaxLength:      128,
		RequireUpper:   true,
		RequireLower:   true,
		RequireNumber:  true,
		RequireSpecial: false,
	}
}

/*
 * Validate checks if password meets policy requirements.
 */
func (p *PasswordPolicy) Validate(password string) error {
	length := len(password)

	if length < p.MinLength {
		return fmt.Errorf("password must be at least %d characters", p.MinLength)
	}

	if length > p.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", p.MaxLength)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if p.RequireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if p.RequireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if p.RequireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if p.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}
