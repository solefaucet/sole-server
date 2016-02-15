package v1

import (
	"net/http"
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

		c.JSON(http.StatusOK, income)
	}
}

// RewardList returns user's reward list as response
func RewardList(
	getRewardIncomesSince dependencyGetRewardIncomesSince,
	getRewardIncomesUntil dependencyGetRewardIncomesUntil,
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
		result := []models.Income{}
		var syserr *errors.Error
		if isSince {
			result, syserr = getRewardIncomesSince(authToken.UserID, t, limit)
		} else {
			result, syserr = getRewardIncomesUntil(authToken.UserID, t, limit)
		}

		// response with result or error
		if syserr != nil {
			c.AbortWithError(http.StatusInternalServerError, syserr)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
