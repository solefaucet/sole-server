package v1

import (
	"net/http"
	"time"

	"github.com/freeusd/solebtc/models"
	"github.com/gin-gonic/gin"
)

// WithdrawalList returns user's withdrawal list as response
func WithdrawalList(
	getWithdrawalsSince dependencyGetWithdrawalsSince,
	getWithdrawalsUntil dependencyGetWithdrawalsUntil,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse pagination args
		isSince, separator, limit, err := parsePagination(c)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// get result according to args
		t := time.Unix(separator, 0)
		result := []models.Withdrawal{}
		if isSince {
			result, err = getWithdrawalsSince(authToken.UserID, t, limit)
		} else {
			result, err = getWithdrawalsUntil(authToken.UserID, t, limit)
		}

		// response with result or error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
