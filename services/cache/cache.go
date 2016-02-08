package cache

import "time"

// Cache defines interface that one should implement
type Cache interface {
	GetBitcoinPrice() int64

	GetTotalRewardToday() int64
	IncrementTotalReward(time.Time, int64)
}
