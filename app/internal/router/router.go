/*
 * GoAstra Backend - HTTP Router
 *
 * Configures Gin router with middleware, routes, and handlers.
 * Provides centralized route registration and middleware chain.
 */
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/channdev/goastra/app/internal/config"
	"github.com/channdev/goastra/app/internal/database"
	"github.com/channdev/goastra/app/internal/logger"
	"github.com/channdev/goastra/app/internal/middleware"
)

/*
 * Router wraps Gin engine with dependencies.
 */
type Router struct {
	engine *gin.Engine
	logger *logger.Logger
	db     *database.DB
	config *config.Config
}

/*
 * New creates a new router instance with all middleware and routes.
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
	/* Recovery middleware for panic handling */
	r.engine.Use(gin.Recovery())

	/* CORS middleware */
	r.engine.Use(middleware.CORS(r.config))

	/* Request logging middleware */
	r.engine.Use(middleware.RequestLogger(r.logger))

	/* Request ID middleware */
	r.engine.Use(middleware.RequestID())
}

func (r *Router) setupRoutes() {
	/* Health check endpoint */
	r.engine.GET("/health", r.healthCheck)

	/* Serve static files in production */
	if r.config.IsProduction() {
		r.engine.Static("/assets", "./public/browser/assets")
		r.engine.StaticFile("/favicon.ico", "./public/browser/favicon.ico")
		r.engine.NoRoute(func(c *gin.Context) {
			c.File("./public/browser/index.html")
		})
	}

	/* API versioning */
	v1 := r.engine.Group("/api/v1")
	{
		/* Public routes */
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.handleLogin)
			auth.POST("/register", r.handleRegister)
			auth.POST("/refresh", r.handleRefresh)
			auth.POST("/logout", r.handleLogout)
		}

		/* Protected routes */
		protected := v1.Group("")
		protected.Use(middleware.Auth(r.config.JWTSecret))
		{
			/* User management */
			users := protected.Group("/users")
			{
				users.GET("", r.handleListUsers)
				users.GET("/:id", r.handleGetUser)
				users.PUT("/:id", r.handleUpdateUser)
				users.DELETE("/:id", r.handleDeleteUser)
			}

			/* Profile endpoint */
			protected.GET("/profile", r.handleGetProfile)
			protected.PUT("/profile", r.handleUpdateProfile)
		}
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

/* Placeholder handlers - to be implemented in handlers package */

func (r *Router) handleLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleRegister(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleRefresh(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleLogout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleListUsers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleGetUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleUpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleDeleteUser(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleGetProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (r *Router) handleUpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
