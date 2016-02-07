package middlewares

import (
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/contrib/cors"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// CORS allow cross domain resources sharing
func CORS() gin.HandlerFunc {
	config := cors.Config{}
	config.AllowedHeaders = []string{"*"}
	config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	config.AbortOnError = true
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.MaxAge = time.Hour * 12
	return cors.New(config)
}
