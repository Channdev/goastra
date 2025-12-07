/*
 * GoAstra Backend - Authentication Middleware
 *
 * JWT token validation and user context injection.
 * Protects routes requiring authentication.
 */
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
 * Auth returns middleware that validates JWT tokens.
 * Extracts user information and sets it in the request context.
 */
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
				"code":  "AUTH_HEADER_MISSING",
			})
			c.Abort()
			return
		}

		/* Extract Bearer token */
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, expected 'Bearer <token>'",
				"code":  "AUTH_FORMAT_INVALID",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		/* Parse and validate token */
		claims, err := validateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		/* Set user context for downstream handlers */
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

/*
 * OptionalAuth validates JWT if present but doesn't require it.
 * Useful for endpoints that behave differently for authenticated users.
 */
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := validateToken(tokenString, jwtSecret)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

/*
 * RequireRole returns middleware that restricts access to specific roles.
 */
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "role information not found",
				"code":  "ROLE_MISSING",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions for this action",
			"code":  "PERMISSION_DENIED",
		})
		c.Abort()
	}
}

func validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

/*
 * GetUserID extracts user ID from context.
 */
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

/*
 * GetUserRole extracts user role from context.
 */
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(string), true
}
