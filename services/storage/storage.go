package storage

import (
	"time"

	"github.com/solefaucet/sole-server/models"
)

// Storage defines interface that one should implement
type Storage interface {
	// User
	GetUserByID(int64) (models.User, error)
	GetUserByEmail(string) (models.User, error)
	CreateUser(models.User) error
	UpdateUserStatus(int64, string) error
	GetReferees(userID int64, limit, offset int64) ([]models.User, error)
	GetNumberOfReferees(userID int64) (int64, error)
	GetWithdrawableUsers(minAmount float64) ([]models.User, error)

	// AuthToken
	GetAuthToken(string) (models.AuthToken, error)
	CreateAuthToken(models.AuthToken) error
	DeleteAuthToken(string) error

	// Session
	GetSessionByToken(string) (models.Session, error)
	UpsertSession(models.Session) error

	// TotalReward
	GetLatestTotalReward() (models.TotalReward, error)

	// RewardRate
	GetRewardRatesByType(string) ([]models.RewardRate, error)

	// Config
	GetLatestConfig() (models.Config, error)

	// Income
	CreateRewardIncome(models.Income, time.Time) error
	GetRewardIncomes(userID int64, limit, offset int64) ([]models.Income, error)
	GetNumberOfRewardIncomes(userID int64) (int64, error)

	// Withdrawal
	CreateWithdrawal(models.Withdrawal) error
	GetWithdrawals(userID int64, limit, offset int64) ([]models.Withdrawal, error)
	GetNumberOfWithdrawals(userID int64) (int64, error)
	GetPendingWithdrawals() ([]models.Withdrawal, error)
	UpdateWithdrawalStatusToProcessing(ids []int64) error
	UpdateWithdrawalStatusToProcessed(ids []int64, transactionID string) error

	// Offerwow
	GetNumberOfOfferwowEvents(eventID string) (int64, error)
	CreateOfferwowIncome(income models.Income, eventID string) error

	// Superrewards
	GetNumberOfSuperrewardsOffers(transactionID string, userID int64) (int64, error)
	CreateSuperrewardsIncome(income models.Income, transactionID, offerID string) error

	// Clixwall
	GetNumberOfClixwallOffers(offerID string, userID int64) (int64, error)
	CreateClixwallIncome(income models.Income, offerID string) error

	// Ptcwall
	CreatePtcwallIncome(income models.Income) error
}
