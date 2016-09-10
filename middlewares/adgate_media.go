package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdgateMediaAuthRequired rejects request if client ip is not in the list
func AdgateMediaAuthRequired(whitelistIPs string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range strings.Split(whitelistIPs, ",") {
		ips[v] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := ips[c.ClientIP()]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if c.Query("status") != "1" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
