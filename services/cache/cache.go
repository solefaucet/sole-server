package cache

import (
	"time"

	"github.com/freeusd/solebtc/models"
)

// Cache defines interface that one should implement
type Cache interface {
	GetBitcoinPrice() int64

	GetLatestTotalReward() int64
	IncrementTotalReward(time.Time, int64)

	GetRewardRatesByType(string) []models.RewardRate
	SetRewardRates(string, []models.RewardRate)

	GetLatestConfig() models.Config
	SetLatestConfig(models.Config)
}
