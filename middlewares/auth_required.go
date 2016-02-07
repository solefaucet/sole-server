package middlewares

import (
	"net/http"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

type authRequiredDependencyGetAuthToken func(authTokenString string) (models.AuthToken, *errors.Error)

// AuthRequired checks if user is authorized
func AuthRequired(
	getAuthToken authRequiredDependencyGetAuthToken,
	authTokenLifetime time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := getAuthToken(c.Request.Header.Get("Auth-Token"))

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
