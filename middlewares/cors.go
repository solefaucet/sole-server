package middlewares

import (
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS allow cross domain resources sharing
func CORS() gin.HandlerFunc {
	config := cors.Config{}
	config.AllowedHeaders = []string{"Content-Type", "Auth-Token", "X-Geetest-Challenge", "X-Geetest-Validate", "X-Geetest-Seccode"}
	config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}
	config.AbortOnError = true
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.MaxAge = time.Hour * 12
	return cors.New(config)
}
