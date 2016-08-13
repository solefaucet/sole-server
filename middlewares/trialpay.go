package middlewares

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/models"
)

// TrialpayAuthRequired rejects request if client ip is not in the list and signature not match
func TrialpayAuthRequired(whitelistIPs, notificationKey string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range strings.Split(whitelistIPs, ",") {
		ips[v] = struct{}{}
	}

	// hardcode it first
	format := "70.42.249.%d"
	for i := 0; i <= 255; i++ {
		ips[fmt.Sprintf(format, i)] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := ips[c.ClientIP()]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		signature := c.Request.Header.Get("TrialPay-HMAC-MD5")
		h := hmac.New(md5.New, []byte(notificationKey))
		h.Write([]byte(c.Request.URL.RawQuery))
		if sign := fmt.Sprintf("%x", h.Sum(nil)); sign != signature {
			httprequest, _ := httputil.DumpRequest(c.Request, true)
			logrus.WithFields(logrus.Fields{
				"event":       models.EventTrialpayInvalidSignature,
				"user_id":     c.Query("sid"),
				"order_id":    c.Query("oid"),
				"amount":      c.Query("reward_amount"),
				"signature":   sign,
				"q_signature": signature,
				"request":     string(httprequest),
			}).Error("signature not matched")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
