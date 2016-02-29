package v1

import (
	"encoding/json"
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
	insertIncome dependencyInsertIncome,
	broadcast dependencyBroadcast,
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
			c.AbortWithStatus(statusCodeTooManyRequests)
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
		rewardReferer := reward * getSystemConfig().RefererRewardRate / 100

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

		// cache delta income
		deltaIncome := struct {
			BitcoinAddress string    `json:"bitcoin_address"`
			Amount         float64   `json:"amount"`
			Type           string    `json:"type"`
			Time           time.Time `json:"time"`
		}{user.BitcoinAddress, utils.HumanReadableBTC(reward), "reward", now}
		insertIncome(deltaIncome)

		// broadcast delta income to all clients
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

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

		// response
		getRewards(c, authToken.UserID, getRewardIncomesSince, getRewardIncomesUntil)
	}
}

// RefereeRewardList returns user's referee's reward list as response
func RefereeRewardList(
	getUserByID dependencyGetUserByID,
	getRewardIncomesSince dependencyGetRewardIncomesSince,
	getRewardIncomesUntil dependencyGetRewardIncomesUntil,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// check if user is referer of :referee_id
		refereeID, _ := strconv.ParseInt(c.Param("referee_id"), 10, 64)
		referee, _ := getUserByID(refereeID)
		if referee.HasReferer() && referee.RefererID != authToken.UserID {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// response
		getRewards(c, refereeID, getRewardIncomesSince, getRewardIncomesUntil)
	}
}

// common get rewards logic
func getRewards(
	c *gin.Context,
	userID int64,
	getRewardIncomesSince dependencyGetRewardIncomesSince,
	getRewardIncomesUntil dependencyGetRewardIncomesUntil,
) {
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
		result, syserr = getRewardIncomesSince(userID, t, limit)
	} else {
		result, syserr = getRewardIncomesUntil(userID, t, limit)
	}

	// response with result or error
	if syserr != nil {
		c.AbortWithError(http.StatusInternalServerError, syserr)
		return
	}

	c.JSON(http.StatusOK, result)
}
