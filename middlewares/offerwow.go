package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OfferwowAuthRequired rejects request if offerwow key is empty
func OfferwowAuthRequired(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
