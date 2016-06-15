package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OfferwowEnableRequired rejects request if offerwow is disabled
func OfferwowEnableRequired(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
