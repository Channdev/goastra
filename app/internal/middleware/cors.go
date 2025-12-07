/*
 * GoAstra Backend - CORS Middleware
 *
 * Handles Cross-Origin Resource Sharing headers for API access.
 * Configurable origins, methods, and headers from config.
 */
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/channdev/goastra/app/internal/config"
)

/*
 * CORS returns middleware that handles CORS headers.
 * Reads configuration from the provided Config struct.
 */
func CORS(cfg *config.Config) gin.HandlerFunc {
	allowedMethods := strings.Join(cfg.CORSAllowedMethods, ", ")
	allowedHeaders := strings.Join(cfg.CORSAllowedHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		/* Check if origin is allowed */
		if isOriginAllowed(origin, cfg.CORSAllowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if contains(cfg.CORSAllowedOrigins, "*") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		/* Handle preflight requests */
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}

	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}

	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
