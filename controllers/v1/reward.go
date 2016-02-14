package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/utils"
)

// GetReward randomly gives users reward
func GetReward(
	getUserByID dependencyGetUserByID,
	getLatestTotalReward dependencyGetLatestTotalReward,
	getSystemConfig dependencyGetSystemConfig,
	getRewardRatesByType dependencyGetRewardRatesByType,
	createRewardIncome dependencyCreateRewardIncome,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)
		now := time.Now()

		// get user
		user, err := getUserByID(authToken.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// check last rewarded time
		if user.RewardedAt.Add(time.Second * time.Duration(user.RewardInterval)).After(now) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// get random bitcoin reward in Satonish
		latestTotalReward := getLatestTotalReward()
		rewardRateType := models.RewardRateTypeLess
		if latestTotalReward.IsSameDay(now) && latestTotalReward.Total > getSystemConfig().TotalRewardThreshold {
			rewardRateType = models.RewardRateTypeMore
		}
		rewardRates := getRewardRatesByType(rewardRateType)
		rewardUSD := utils.RandomReward(rewardRates)
		btcPrice := getSystemConfig().BitcoinPrice
		reward := utils.MachineReadableBTC(float64(rewardUSD) / float64(btcPrice))
		rewardReferer := int64(float64(reward) * getSystemConfig().RefererRewardRate)

		// create income reward
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeReward,
			Income:        reward,
			RefererIncome: rewardReferer,
		}
		if err := createRewardIncome(income, now); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// RewardList returns user's reward list as response
func RewardList(
	getRewardIncomesSince dependencyGetRewardIncomesSince,
	getRewardIncomesUntil dependencyGetRewardIncomesUntil,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse limit
		queryLimit := c.DefaultQuery("limit", "10")
		limit, err := strconv.ParseInt(queryLimit, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		limit = max(limit, 100)

		// parse time
		querySince := c.Query("since")
		queryUntil := c.Query("until")
		var f func(int64, time.Time, int64) ([]models.Income, *errors.Error)
		var t time.Time
		switch {
		case querySince != "":
			f = getRewardIncomesSince
			t, err = time.Parse(time.RFC3339, querySince)
		case queryUntil != "":
			f = getRewardIncomesUntil
			t, err = time.Parse(time.RFC3339, queryUntil)
		default:
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// response with incomes
		incomes, syserr := f(authToken.UserID, t, limit)
		if syserr != nil {
			c.AbortWithError(http.StatusInternalServerError, syserr)
			return
		}

		c.JSON(http.StatusOK, incomes)
	}
}
