package v1

import (
	"net/http"
	"time"

	"github.com/solefaucet/solebtc/models"
	"github.com/gin-gonic/gin"
)

// WithdrawalList returns user's withdrawal list as response
func WithdrawalList(
	getWithdrawals dependencyGetWithdrawals,
	getNumberOfWithdrawals dependencyGetNumberOfWithdrawals,
	constructTxURL dependencyConstructTxURL,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse pagination args
		limit, offset, err := parsePagination(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		withdrawals, err := getWithdrawals(authToken.UserID, limit, offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfWithdrawals(authToken.UserID)
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

		c.JSON(http.StatusOK, paginationResult(result, count))
	}
}
