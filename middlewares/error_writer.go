package middlewares

import "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"

// ErrorWriter writes last error into response body if not written yet
func ErrorWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		if err := c.Errors.ByType(gin.ErrorTypeAny).Last(); err != nil {
			c.JSON(-1, err)
		}
	}
}
