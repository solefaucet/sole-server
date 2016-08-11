package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/models"
)

// OfferwallList returns user's offerwall list as response
func OfferwallList(
	getOfferwallIncomes dependencyGetOfferwallIncomes,
	getNumberOfOfferwallIncomes dependencyGetNumberOfOfferwallIncomes,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse pagination args
		limit, offset, err := parsePagination(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		offerwalls, err := getOfferwallIncomes(authToken.UserID, limit, offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfOfferwallIncomes(authToken.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult(offerwalls, count))
	}
}
