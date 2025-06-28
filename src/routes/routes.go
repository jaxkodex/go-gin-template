package routes

import "github.com/gin-gonic/gin"

// RegisterRoutes adds all route modules to the Gin router
func RegisterRoutes(r *gin.Engine) {
	HealthRoute(r)
	// Add more route registrations here as your app grows
}
