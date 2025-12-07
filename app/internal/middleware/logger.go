/*
 * GoAstra Backend - Request Logger Middleware
 *
 * Structured logging for HTTP requests with timing information.
 * Captures request metadata for debugging and monitoring.
 */
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/channdev/goastra/app/internal/logger"
)

/*
 * RequestLogger returns middleware that logs HTTP requests.
 * Captures method, path, status, latency, and client information.
 */
func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		/* Process request */
		c.Next()

		/* Calculate latency */
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		/* Build log fields */
		fields := []interface{}{
			"method", c.Request.Method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"latency_ms", latency.Milliseconds(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		if query != "" {
			fields = append(fields, "query", query)
		}

		if requestID, exists := c.Get("request_id"); exists {
			fields = append(fields, "request_id", requestID)
		}

		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, "user_id", userID)
		}

		/* Log based on status code */
		if statusCode >= 500 {
			log.Errorw("request completed with server error", fields...)
		} else if statusCode >= 400 {
			log.Warnw("request completed with client error", fields...)
		} else {
			log.Infow("request completed", fields...)
		}
	}
}

/*
 * RequestID returns middleware that generates unique request IDs.
 * Uses existing X-Request-ID header or generates a new one.
 */
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func generateRequestID() string {
	/* Simple timestamp-based ID */
	return time.Now().Format("20060102150405.000000")
}
