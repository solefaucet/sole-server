package middlewares

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Logger returns a middleware that logs all request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// mark start time
		start := time.Now()

		// process request
		c.Next()

		// log after processing
		status := c.Writer.Status()
		fields := logrus.Fields{
			"method":           c.Request.Method,
			"path":             c.Request.URL.Path,
			"query":            c.Request.URL.Query().Encode(),
			"http_status_code": status,
			"response_time":    time.Since(start),
			"ip":               c.ClientIP(),
		}
		if err := c.Errors.ByType(gin.ErrorTypeAny).Last(); err != nil {
			fields["error"] = err
		}

		entry := logrus.WithFields(fields)
		if status == http.StatusInternalServerError {
			entry.Error("http error")
			return
		}

		entry.Info("http access")
	}
}
