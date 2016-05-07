package middlewares

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const k = 1 << 10

// RecoveryWithWriter returns a middleware that
// recovers from any panics and writes a 500 response code
func RecoveryWithWriter(out io.Writer) gin.HandlerFunc {
	logger := log.New(out, "", log.LstdFlags)

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, true)
				buf := make([]byte, 4*k)
				n := runtime.Stack(buf, false)

				logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s\n", httprequest, err, buf[:n])
				logrus.WithFields(logrus.Fields{
					"request": string(httprequest),
					"error":   err,
					"stack":   string(buf[:n]),
				}).Error("panic")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
