/*
 * GoAstra Backend - JWT Authentication
 *
 * Token generation and validation for user authentication.
 * Supports access tokens and refresh tokens.
 */
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
 * Claims represents the JWT token payload structure.
 */
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

/*
 * RefreshClaims represents refresh token payload.
 */
type RefreshClaims struct {
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

/*
 * TokenPair contains both access and refresh tokens.
 */
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

/*
 * JWTManager handles token generation and validation.
 */
type JWTManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

/*
 * NewJWTManager creates a new JWT manager instance.
 */
func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        "goastra",
	}
}

/*
 * GenerateTokenPair creates both access and refresh tokens.
 */
func (m *JWTManager) GenerateTokenPair(userID uint, email, role string) (*TokenPair, error) {
	accessToken, err := m.generateAccessToken(userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := m.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (m *JWTManager) generateAccessToken(userID uint, email, role string) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *JWTManager) generateRefreshToken(userID uint) (string, error) {
	now := time.Now()

	claims := &RefreshClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

/*
 * ValidateAccessToken verifies and extracts claims from an access token.
 */
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

/*
 * ValidateRefreshToken verifies and extracts claims from a refresh token.
 */
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	return claims, nil
}

/*
 * RefreshTokens generates a new token pair using a valid refresh token.
 */
func (m *JWTManager) RefreshTokens(refreshToken string, email, role string) (*TokenPair, error) {
	claims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return m.GenerateTokenPair(claims.UserID, email, role)
}
