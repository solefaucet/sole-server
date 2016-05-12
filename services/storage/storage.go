package storage

import (
	"time"

	"github.com/freeusd/solebtc/models"
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
	GetWithdrawableUsers() ([]models.User, error)

	// AuthToken
	GetAuthToken(string) (models.AuthToken, error)
	CreateAuthToken(models.AuthToken) error
	DeleteAuthToken(string) error

	// Session
	GetSessionByToken(string) (models.Session, error)
	UpsertSession(models.Session) error

	// TotalReward
	IncrementTotalReward(time.Time, float64) error
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
	UpdateWithdrawalStatusToProcessing(id int64) error
	UpdateWithdrawalStatusToProcessed(id int64, transactionID string) error
}
