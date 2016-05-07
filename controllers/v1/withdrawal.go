package v1

import (
	"net/http"
	"time"

	"github.com/freeusd/solebtc/errors"
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
		var syserr *errors.Error
		if isSince {
			result, syserr = getWithdrawalsSince(authToken.UserID, t, limit)
		} else {
			result, syserr = getWithdrawalsUntil(authToken.UserID, t, limit)
		}

		// response with result or error
		if syserr != nil {
			c.AbortWithError(http.StatusInternalServerError, syserr)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
