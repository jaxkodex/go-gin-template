package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/middleware"
)

// RegisterRoutes adds all route modules to the Gin router.
func RegisterRoutes(r *gin.Engine, auth *middleware.AuthMiddleware) {
	// Public routes
	HealthRoute(r)

	// Protected group — requires authentication when auth is enabled.
	protected := r.Group("/")
	protected.Use(auth.Authenticate())
	{
		// Register protected routes here as your app grows.
		_ = protected
	}
}
