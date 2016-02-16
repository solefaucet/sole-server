package cache

import (
	"time"

	"github.com/freeusd/solebtc/models"
)

// Cache defines interface that one should implement
type Cache interface {
	GetLatestTotalReward() models.TotalReward
	IncrementTotalReward(time.Time, int64)

	GetRewardRatesByType(string) []models.RewardRate
	SetRewardRates(string, []models.RewardRate)

	GetLatestConfig() models.Config
	SetLatestConfig(models.Config)
	UpdateBitcoinPrice(int64)
}
