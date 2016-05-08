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
	GetRefereesSince(userID, id, limit int64) ([]models.User, error)
	GetRefereesUntil(userID, id, limit int64) ([]models.User, error)
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
	GetRewardIncomesSince(userID int64, since time.Time, limit int64) ([]models.Income, error)
	GetRewardIncomesUntil(userID int64, until time.Time, limit int64) ([]models.Income, error)

	// Withdrawal
	CreateWithdrawal(models.Withdrawal) error
	GetWithdrawalsSince(userID int64, since time.Time, limit int64) ([]models.Withdrawal, error)
	GetWithdrawalsUntil(userID int64, until time.Time, limit int64) ([]models.Withdrawal, error)
	GetUnprocessedWithdrawals() ([]models.Withdrawal, error)
}
