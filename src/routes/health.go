package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/database"
)

// HealthRoute registers the /health endpoint
func HealthRoute(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		if err := database.DB.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
