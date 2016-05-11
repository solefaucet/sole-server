package v1

import (
	"time"

	"github.com/gorilla/websocket"

	"github.com/freeusd/solebtc/models"
)

// dependencies
type (
	// user
	dependencyGetUserByID      func(int64) (models.User, error)
	dependencyGetUserByEmail   func(string) (models.User, error)
	dependencyCreateUser       func(models.User) error
	dependencyUpdateUserStatus func(int64, string) error
	dependencyGetRefereesSince func(userID int64, sinceID int64, limit int64) ([]models.User, error)
	dependencyGetRefereesUntil func(userID int64, untilID int64, limit int64) ([]models.User, error)

	// auth token
	dependencyCreateAuthToken func(models.AuthToken) error
	dependencyDeleteAuthToken func(string) error

	// session
	dependencyUpsertSession     func(models.Session) error
	dependencyGetSessionByToken func(string) (models.Session, error)

	// email
	dependencySendEmail func(recipients []string, subject string, html string) error

	// total reward
	dependencyGetLatestTotalReward func() models.TotalReward
	dependencyIncrementTotalReward func(time.Time, int64) error

	// reward rate
	dependencyGetRewardRatesByType func(string) []models.RewardRate

	// system config
	dependencyGetSystemConfig func() models.Config

	// income
	dependencyCreateRewardIncome    func(models.Income, time.Time) error
	dependencyGetRewardIncomesSince func(userID int64, since time.Time, limit int64) ([]models.Income, error)
	dependencyGetRewardIncomesUntil func(userID int64, until time.Time, limit int64) ([]models.Income, error)
	dependencyInsertIncome          func(interface{}) // cache for broadcasting

	// websocket
	dependencyPutConn          func(*websocket.Conn)
	dependencyBroadcast        func([]byte)
	dependencyGetUsersOnline   func() int
	dependencyGetLatestIncomes func() []interface{}

	// withdrawals
	dependencyGetWithdrawalsSince func(userID int64, since time.Time, limit int64) ([]models.Withdrawal, error)
	dependencyGetWithdrawalsUntil func(userID int64, until time.Time, limit int64) ([]models.Withdrawal, error)
	dependencyConstructTxURL      func(tx string) string

	// validation
	dependencyValidateAddress func(string) (bool, error)
)
