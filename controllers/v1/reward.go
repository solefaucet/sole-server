package v1

import (
	"fmt"
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
	getBitcoinPrice dependencyGetBitcoinPrice,
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
		btcPrice := getBitcoinPrice()
		reward := utils.MachineReadableBTC(float64(rewardUSD) / float64(btcPrice))
		rewardReferer := int64(float64(reward) * getSystemConfig().RefererRewardRate)

		// create income reward
		if err := createRewardIncome(user.ID, user.RefererID, reward, rewardReferer, now); err != nil {
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

		// parse timestamp
		queryTimestamp := c.Query("timestamp")
		timestamp, err := strconv.ParseInt(queryTimestamp, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		t := time.Unix(timestamp, 0)

		// parse since or until
		incomes := []models.Income{}
		var syserr *errors.Error
		switch c.Query("type") {
		case "since":
			incomes, syserr = getRewardIncomesSince(authToken.UserID, t, limit)
		case "until":
			incomes, syserr = getRewardIncomesUntil(authToken.UserID, t, limit)
		default:
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if syserr != nil {
			fmt.Printf("got error: %v\n", syserr)
			c.AbortWithError(http.StatusInternalServerError, syserr)
			return
		}

		c.JSON(http.StatusOK, incomes)
	}
}
