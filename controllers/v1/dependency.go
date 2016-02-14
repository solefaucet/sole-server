package v1

import (
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// dependencies
type (
	// user
	dependencyGetUserByID    func(int64) (models.User, *errors.Error)
	dependencyGetUserByEmail func(string) (models.User, *errors.Error)
	dependencyCreateUser     func(models.User) *errors.Error
	dependencyUpdateUser     func(models.User) *errors.Error

	// auth token
	dependencyCreateAuthToken func(models.AuthToken) *errors.Error
	dependencyDeleteAuthToken func(string) *errors.Error

	// session
	dependencyUpsertSession     func(models.Session) *errors.Error
	dependencyGetSessionByToken func(string) (models.Session, *errors.Error)

	// email
	dependencySendEmail func(recipients []string, subject string, html string) *errors.Error

	// bitcoin price
	dependencyGetBitcoinPrice func() int64

	// total reward
	dependencyGetLatestTotalReward func() models.TotalReward
	dependencyIncrementTotalReward func(time.Time, int64) *errors.Error

	// reward rate
	dependencyGetRewardRatesByType func(string) []models.RewardRate

	// system config
	dependencyGetSystemConfig func() models.Config

	// income
	dependencyCreateRewardIncome    func(userID, refererID, reward, rewardReferer int64, now time.Time) *errors.Error
	dependencyGetRewardIncomesSince func(userID int64, since time.Time, limit int64) ([]models.Income, *errors.Error)
	dependencyGetRewardIncomesUntil func(userID int64, until time.Time, limit int64) ([]models.Income, *errors.Error)
)
