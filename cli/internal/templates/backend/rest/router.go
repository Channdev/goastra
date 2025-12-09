/*
 * GoAstra CLI - REST Router Template
 *
 * Generates the router configuration for REST APIs.
 * Sets up middleware chain, routes, and handlers.
 */
package rest

// RouterGo returns the router.go template for REST APIs.
func RouterGo() string {
	return `package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"app/internal/config"
	"app/internal/database"
	"app/internal/logger"
	"app/internal/middleware"
)

/*
 * Router wraps Gin engine with application dependencies.
 */
type Router struct {
	engine *gin.Engine
	logger *logger.Logger
	db     *database.DB
	config *config.Config
}

/*
 * New creates a new router instance with middleware and routes configured.
 */
func New(log *logger.Logger, db *database.DB, cfg *config.Config) *Router {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := &Router{
		engine: gin.New(),
		logger: log,
		db:     db,
		config: cfg,
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

/*
 * Handler returns the underlying HTTP handler.
 */
func (r *Router) Handler() http.Handler {
	return r.engine
}

func (r *Router) setupMiddleware() {
	// Recovery middleware
	r.engine.Use(gin.Recovery())

	// CORS middleware
	r.engine.Use(middleware.CORS(r.config.CORSOrigins))

	// Request ID middleware
	r.engine.Use(middleware.RequestID())

	// Request logging middleware
	r.engine.Use(r.requestLogger())
}

func (r *Router) setupRoutes() {
	// Health check endpoint
	r.engine.GET("/health", r.healthCheck)

	// API versioning
	v1 := r.engine.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.handleLogin)
			auth.POST("/register", r.handleRegister)
			auth.POST("/refresh", r.handleRefresh)
			auth.POST("/logout", r.handleLogout)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(r.config.JWTSecret))
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("", r.handleListUsers)
				users.GET("/:id", r.handleGetUser)
				users.PUT("/:id", r.handleUpdateUser)
				users.DELETE("/:id", r.handleDeleteUser)
			}

			// Profile endpoints
			protected.GET("/profile", r.handleGetProfile)
			protected.PUT("/profile", r.handleUpdateProfile)
		}
	}

	// Serve static files in production
	if r.config.IsProduction() {
		r.setupStaticFiles()
	}
}

func (r *Router) setupStaticFiles() {
	staticDir := "./public/browser"

	// Serve assets
	r.engine.Static("/assets", staticDir+"/assets")
	r.engine.StaticFile("/favicon.ico", staticDir+"/favicon.ico")

	// SPA fallback
	r.engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File(staticDir + "/index.html")
	})
}

/*
 * requestLogger returns middleware that logs HTTP requests with style.
 */
func (r *Router) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		// Skip OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			return
		}

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		// Get handler file from route
		handler := getHandlerFile(path)

		// Format and print log
		fmt.Fprintf(os.Stdout, "%s | %s%3d%s | %s%8s%s | %s%-12s%s %s%-7s%s \"%s\"\n",
			time.Now().Format("15:04:05"),
			statusColor(status), status, colorReset,
			colorCyan, formatLatency(latency), colorReset,
			colorMagenta, handler, colorReset,
			methodColor(method), method, colorReset,
			path,
		)
	}
}

/*
 * healthCheck returns service health status.
 */
func (r *Router) healthCheck(c *gin.Context) {
	status := "healthy"
	dbStatus := "connected"

	if r.db != nil {
		if err := r.db.Health(); err != nil {
			dbStatus = "disconnected"
			status = "degraded"
		}
	} else {
		dbStatus = "not configured"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"database": dbStatus,
		"version":  "1.0.0",
	})
}

// Placeholder handlers - implement in handlers package
func (r *Router) handleLogin(c *gin.Context)         { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleRegister(c *gin.Context)      { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleRefresh(c *gin.Context)       { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleLogout(c *gin.Context)        { c.JSON(200, gin.H{"message": "logged out"}) }
func (r *Router) handleListUsers(c *gin.Context)     { c.JSON(200, gin.H{"data": []interface{}{}, "total": 0}) }
func (r *Router) handleGetUser(c *gin.Context)       { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleUpdateUser(c *gin.Context)    { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleDeleteUser(c *gin.Context)    { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleGetProfile(c *gin.Context)    { c.JSON(501, gin.H{"error": "not implemented"}) }
func (r *Router) handleUpdateProfile(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }

// Logging helpers
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

var routeHandlers = map[string]string{
	"/api/v1/auth":    "auth.go",
	"/api/v1/users":   "users.go",
	"/api/v1/profile": "profile.go",
	"/health":         "router.go",
}

func getHandlerFile(path string) string {
	for route, handler := range routeHandlers {
		if strings.HasPrefix(path, route) {
			return handler
		}
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		return parts[2] + ".go"
	}
	return "handler.go"
}

func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fµs", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func statusColor(code int) string {
	switch {
	case code >= 500:
		return colorRed
	case code >= 400:
		return colorYellow
	case code >= 300:
		return colorCyan
	default:
		return colorGreen
	}
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorGreen
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	default:
		return colorCyan
	}
}
`
}
