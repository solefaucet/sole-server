package middlewares

import (
	"net/http"
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/gin-gonic/gin"
)

type authRequiredDependencyGetAuthToken func(authTokenString string) (models.AuthToken, *errors.Error)

// AuthRequired checks if user is authorized
func AuthRequired(
	getAuthToken authRequiredDependencyGetAuthToken,
	authTokenLifetime time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authTokenHeader := c.Request.Header.Get("Auth-Token")
		if authTokenHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken, err := getAuthToken(authTokenHeader)
		if err != nil && err.ErrCode != errors.ErrCodeNotFound {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if authToken.CreatedAt.Add(authTokenLifetime).Before(time.Now()) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("auth_token", authToken)
		c.Next()
	}
}
