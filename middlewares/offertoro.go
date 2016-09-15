package middlewares

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/models"
)

// OffertoroAuthRequired rejects request if signature not match
func OffertoroAuthRequired(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := fmt.Sprintf("%v-%v-%v", c.Query("oid"), c.Query("user_id"), secretKey)
		if sign := fmt.Sprintf("%x", md5.Sum([]byte(data))); sign != c.Query("sig") {
			httprequest, _ := httputil.DumpRequest(c.Request, true)
			logrus.WithFields(logrus.Fields{
				"event":       models.EventOffertoroInvalidSignature,
				"user_id":     c.Query("user_id"),
				"signature":   sign,
				"q_signature": c.Query("sig"),
				"request":     string(httprequest),
			}).Error("signature not matched")
			c.String(http.StatusForbidden, "0")
			return
		}

		c.Next()
	}
}
