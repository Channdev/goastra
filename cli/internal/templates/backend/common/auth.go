/*
 * GoAstra CLI - Auth Template
 *
 * Generates JWT authentication utilities including token
 * generation, validation, and password hashing.
 */
package common

// AuthGo returns the auth.go template with JWT and password utilities.
func AuthGo() string {
	return `package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

/*
 * Claims represents JWT token claims with user information.
 */
type Claims struct {
	UserID uint   ` + "`json:\"user_id\"`" + `
	Email  string ` + "`json:\"email\"`" + `
	Role   string ` + "`json:\"role\"`" + `
	jwt.RegisteredClaims
}

/*
 * TokenPair holds access and refresh tokens.
 */
type TokenPair struct {
	AccessToken  string ` + "`json:\"access_token\"`" + `
	RefreshToken string ` + "`json:\"refresh_token\"`" + `
	ExpiresAt    int64  ` + "`json:\"expires_at\"`" + `
}

/*
 * GenerateToken creates a new JWT token for the given user.
 */
func GenerateToken(userID uint, email, role, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

/*
 * GenerateTokenPair creates both access and refresh tokens.
 */
func GenerateTokenPair(userID uint, email, role, secret string) (*TokenPair, error) {
	accessExpiry := 24 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour

	accessToken, err := GenerateToken(userID, email, role, secret, accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateToken(userID, email, role, secret+"_refresh", refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessExpiry).Unix(),
	}, nil
}

/*
 * ValidateToken parses and validates a JWT token string.
 * Returns the claims if valid, error otherwise.
 */
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

/*
 * HashPassword generates a bcrypt hash of the given password.
 */
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

/*
 * CheckPassword compares a password with its hash.
 * Returns true if they match.
 */
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
`
}
