package storage

import (
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// Storage defines interface that one should implement
type Storage interface {
	// User
	GetUserByID(int64) (models.User, *errors.Error)
	GetUserByEmail(string) (models.User, *errors.Error)
	CreateUser(models.User) *errors.Error
	UpdateUser(models.User) *errors.Error
	IncrementUserBalance(int64, int64, *time.Time) *errors.Error

	// AuthToken
	GetAuthToken(string) (models.AuthToken, *errors.Error)
	CreateAuthToken(models.AuthToken) *errors.Error
	DeleteAuthToken(string) *errors.Error

	// Session
	GetSessionByToken(string) (models.Session, *errors.Error)
	UpsertSession(models.Session) *errors.Error

	// TotalReward
	IncrementTotalReward(time.Time, int64) *errors.Error
	GetLatestTotalReward() (models.TotalReward, *errors.Error)

	// RewardRate
	GetRewardRatesByType(string) ([]models.RewardRate, *errors.Error)
}
