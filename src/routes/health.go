package routes

import (
	"github.com/gin-gonic/gin"
)

// HealthRoute registers the /health endpoint
func HealthRoute(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
