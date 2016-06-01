package cache

import (
	"time"

	"github.com/solefaucet/sole-server/models"
)

// Cache defines interface that one should implement
type Cache interface {
	GetLatestTotalReward() models.TotalReward
	IncrementTotalReward(time.Time, float64)

	GetRewardRatesByType(string) []models.RewardRate
	SetRewardRates(string, []models.RewardRate)

	GetLatestConfig() models.Config
	SetLatestConfig(models.Config)

	InsertIncome(interface{})
	GetLatestIncomes() []interface{}
}
