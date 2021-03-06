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
	GetOfferwallIncomes(userID int64, limit, offset int64) ([]models.Income, error)
	GetNumberOfOfferwallIncomes(userID int64) (int64, error)
	ChargebackIncome(incomeID int64) error

	// Withdrawal
	CreateWithdrawal(models.Withdrawal) error
	GetWithdrawals(userID int64, limit, offset int64) ([]models.Withdrawal, error)
	GetNumberOfWithdrawals(userID int64) (int64, error)
	GetPendingWithdrawals() ([]models.Withdrawal, error)
	UpdateWithdrawalStatusToProcessing(ids []int64) error
	UpdateWithdrawalStatusToProcessed(ids []int64, transactionID string) error

	// Superrewards
	GetNumberOfSuperrewardsOffers(transactionID string, userID int64) (int64, error)
	CreateSuperrewardsIncome(income models.Income, transactionID, offerID string) error

	// Clixwall
	GetNumberOfClixwallOffers(offerID string, userID int64) (int64, error)
	CreateClixwallIncome(income models.Income, offerID string) error

	// Ptcwall
	CreatePtcwallIncome(income models.Income) error

	// Personaly
	GetNumberOfPersonalyOffers(offerID string, userID int64) (int64, error)
	CreatePersonalyIncome(income models.Income, offerID string) error

	// Kiwiwall
	GetNumberOfKiwiwallOffers(transactionID string, userID int64) (int64, error)
	CreateKiwiwallIncome(income models.Income, transactionID, offerID string) error

	// AdscendMedia
	GetAdscendMediaOffer(transactionID string, userID int64) (*models.AdscendMedia, error)
	CreateAdscendMediaIncome(income models.Income, transactionID, offerID string) error

	// AdgateMedia
	GetNumberOfAdgateMediaOffers(transactionID string, userID int64) (int64, error)
	CreateAdgateMediaIncome(income models.Income, transactionID, offerID string) error

	// Offertoro
	GetNumberOfOffertoroOffers(transactionID string, userID int64) (int64, error)
	CreateOffertoroIncome(income models.Income, transactionID, offerID string) error
}
