/*
 * GoAstra Backend - Request Logger Middleware
 *
 * Stylish console logging for HTTP requests with timing and handler info.
 * Captures request metadata for debugging and monitoring.
 */
package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/channdev/goastra/app/internal/logger"
)

/* ANSI color codes */
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
)

/* Route to handler file mapping */
var routeHandlers = map[string]string{
	"/api/v1/auth/login":    "auth.go",
	"/api/v1/auth/register": "auth.go",
	"/api/v1/auth/refresh":  "auth.go",
	"/api/v1/auth/logout":   "auth.go",
	"/api/v1/users":         "users.go",
	"/api/v1/profile":       "profile.go",
	"/health":               "router.go",
}

/*
 * RequestLogger returns middleware that logs HTTP requests with style.
 * Format: timestamp | status | latency | handler | METHOD "path"
 */
func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		/* Process request */
		c.Next()

		/* Calculate latency */
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method

		/* Skip OPTIONS preflight logging for cleaner output */
		if method == "OPTIONS" {
			return
		}

		/* Get handler file from route */
		handler := getHandlerFile(path)

		/* Format latency */
		latencyStr := formatLatency(latency)

		/* Get status color */
		statusColor := getStatusColor(statusCode)

		/* Get method color */
		methodColor := getMethodColor(method)

		/* Format timestamp */
		timestamp := time.Now().Format("15:04:05")

		/* Build stylish log line */
		fmt.Fprintf(os.Stdout, "%s%s%s %s|%s %s%3d%s %s|%s %s%8s%s %s|%s %s%-12s%s %s%s%-7s%s %s\"%s\"%s\n",
			colorGray, timestamp, colorReset,
			colorGray, colorReset,
			statusColor, statusCode, colorReset,
			colorGray, colorReset,
			colorCyan, latencyStr, colorReset,
			colorGray, colorReset,
			colorMagenta, handler, colorReset,
			methodColor, colorBold, method, colorReset,
			colorWhite, path, colorReset,
		)

		/* Also log to structured logger for production/JSON output */
		fields := []interface{}{
			"status", statusCode,
			"latency_ms", latency.Milliseconds(),
			"handler", handler,
			"method", method,
			"path", path,
			"ip", c.ClientIP(),
		}

		if statusCode >= 500 {
			log.Errorw("", fields...)
		} else if statusCode >= 400 {
			log.Warnw("", fields...)
		}
	}
}

/*
 * getHandlerFile returns the handler file name for a given route.
 */
func getHandlerFile(path string) string {
	/* Check exact match first */
	if handler, ok := routeHandlers[path]; ok {
		return handler
	}

	/* Check prefix matches for dynamic routes */
	for route, handler := range routeHandlers {
		if strings.HasPrefix(path, route) {
			return handler
		}
	}

	/* Default based on path segments */
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2] + ".go"
	}

	return "handler.go"
}

/*
 * formatLatency formats duration to a readable string.
 */
func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fµs", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

/*
 * getStatusColor returns ANSI color code based on HTTP status.
 */
func getStatusColor(code int) string {
	switch {
	case code >= 500:
		return colorRed + colorBold
	case code >= 400:
		return colorYellow
	case code >= 300:
		return colorCyan
	case code >= 200:
		return colorGreen
	default:
		return colorWhite
	}
}

/*
 * getMethodColor returns ANSI color code based on HTTP method.
 */
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorGreen
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorCyan
	default:
		return colorWhite
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
