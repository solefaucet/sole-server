package v1

import (
	"time"

	"github.com/gorilla/websocket"

	"github.com/solefaucet/sole-server/models"
)

// dependencies
type (
	// user
	dependencyGetUserByID         func(int64) (models.User, error)
	dependencyGetUserByEmail      func(string) (models.User, error)
	dependencyCreateUser          func(models.User) error
	dependencyUpdateUserStatus    func(int64, string) error
	dependencyGetReferees         func(userID int64, limit, offset int64) ([]models.User, error)
	dependencyGetNumberOfReferees func(userID int64) (int64, error)

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
	dependencyCreateRewardIncome       func(models.Income, time.Time) error
	dependencyCreateOfferwowIncome     func(models.Income, string) error
	dependencyCreateSuperrewardsIncome func(income models.Income, transactionID, offerID string) error
	dependencyGetRewardIncomes         func(userID int64, limit, offset int64) ([]models.Income, error)
	dependencyGetNumberOfRewardIncomes func(userID int64) (int64, error)
	dependencyInsertIncome             func(interface{}) // cache for broadcasting

	// websocket
	dependencyPutConn          func(*websocket.Conn)
	dependencyBroadcast        func([]byte)
	dependencyGetUsersOnline   func() int
	dependencyGetLatestIncomes func() []interface{}

	// withdrawals
	dependencyGetWithdrawals         func(userID int64, limit, offset int64) ([]models.Withdrawal, error)
	dependencyGetNumberOfWithdrawals func(userID int64) (int64, error)
	dependencyConstructTxURL         func(tx string) string

	// validation
	dependencyValidateAddress func(string) (bool, error)

	// captcha
	dependencyRegisterCaptcha func() (string, error)
	dependencyGetCaptchaID    func() string

	// offerwow
	dependencyGetNumberOfOfferwowEvent func(eventID string) (int64, error)

	// superrewards
	dependencyGetSuperrewardsOfferByID func(transactionID string, userID int64) (models.SuperrewardsOffer, error)
)
