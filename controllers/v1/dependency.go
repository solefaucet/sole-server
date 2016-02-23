package v1

import (
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// dependencies
type (
	// user
	dependencyGetUserByID      func(int64) (models.User, *errors.Error)
	dependencyGetUserByEmail   func(string) (models.User, *errors.Error)
	dependencyCreateUser       func(models.User) *errors.Error
	dependencyUpdateUserStatus func(int64, string) *errors.Error
	dependencyGetRefereesSince func(userID int64, sinceID int64, limit int64) ([]models.User, *errors.Error)
	dependencyGetRefereesUntil func(userID int64, untilID int64, limit int64) ([]models.User, *errors.Error)

	// auth token
	dependencyCreateAuthToken func(models.AuthToken) *errors.Error
	dependencyDeleteAuthToken func(string) *errors.Error

	// session
	dependencyUpsertSession     func(models.Session) *errors.Error
	dependencyGetSessionByToken func(string) (models.Session, *errors.Error)

	// email
	dependencySendEmail func(recipients []string, subject string, html string) *errors.Error

	// total reward
	dependencyGetLatestTotalReward func() models.TotalReward
	dependencyIncrementTotalReward func(time.Time, int64) *errors.Error

	// reward rate
	dependencyGetRewardRatesByType func(string) []models.RewardRate

	// system config
	dependencyGetSystemConfig func() models.Config

	// income
	dependencyCreateRewardIncome    func(models.Income, time.Time) *errors.Error
	dependencyGetRewardIncomesSince func(userID int64, since time.Time, limit int64) ([]models.Income, *errors.Error)
	dependencyGetRewardIncomesUntil func(userID int64, until time.Time, limit int64) ([]models.Income, *errors.Error)

	// websocket
	dependencyPutConn func(*websocket.Conn)
)
