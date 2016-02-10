package cache

import (
	"time"

	"github.com/freeusd/solebtc/models"
)

// Cache defines interface that one should implement
type Cache interface {
	GetBitcoinPrice() int64

	GetTotalRewardToday() int64
	IncrementTotalReward(time.Time, int64)

	GetRewardRatesByType(string) []models.RewardRate
	SetRewardRates(string, []models.RewardRate)
}
