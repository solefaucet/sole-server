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
	constructTxURL dependencyConstructTxURL,
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
		withdrawals := []models.Withdrawal{}
		if isSince {
			withdrawals, err = getWithdrawalsSince(authToken.UserID, t, limit)
		} else {
			withdrawals, err = getWithdrawalsUntil(authToken.UserID, t, limit)
		}

		// response with result or error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		result := make([]struct {
			UpdatedAt time.Time `json:"updated_at"`
			Amount    float64   `json:"amount"`
			TxURL     string    `json:"tx_url"`
			Status    int64     `json:"status"`
		}, len(withdrawals))
		for i := range withdrawals {
			result[i].UpdatedAt = withdrawals[i].UpdatedAt
			result[i].Amount = withdrawals[i].Amount
			result[i].TxURL = constructTxURL(withdrawals[i].TransactionID)
			result[i].Status = withdrawals[i].Status
		}

		c.JSON(http.StatusOK, result)
	}
}
