package memory

import (
	"sync"
	"time"

	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/services/cache"
)

// Cache implements cache.Cache interface with memory
type Cache struct {
	totalReward      models.TotalReward
	totalRewardMutex sync.RWMutex

	rewardRatesMapping map[string][]models.RewardRate
	rewardRatesMutex   sync.RWMutex

	config      models.Config
	configMutex sync.RWMutex
}

var _ cache.Cache = &Cache{}

// New creates a new in-memory cache
func New() *Cache {
	return &Cache{
		rewardRatesMapping: make(map[string][]models.RewardRate),
	}
}

// GetLatestTotalReward returns total reward of today
func (c *Cache) GetLatestTotalReward() models.TotalReward {
	c.totalRewardMutex.RLock()
	defer c.totalRewardMutex.RUnlock()
	return c.totalReward
}

// IncrementTotalReward increment total reward today by delta if day matches
func (c *Cache) IncrementTotalReward(t time.Time, delta int64) {
	c.totalRewardMutex.Lock()
	defer c.totalRewardMutex.Unlock()

	if c.totalReward.IsSameDay(t) {
		c.totalReward.Total += delta
	} else {
		c.totalReward = models.TotalReward{CreatedAt: t.UTC(), Total: delta}
	}
}

// GetRewardRatesByType returns reward rates by type
func (c *Cache) GetRewardRatesByType(t string) []models.RewardRate {
	c.rewardRatesMutex.RLock()
	defer c.rewardRatesMutex.RUnlock()
	return c.rewardRatesMapping[t]
}

// SetRewardRates sets reward rates with type
func (c *Cache) SetRewardRates(t string, rates []models.RewardRate) {
	c.rewardRatesMutex.Lock()
	defer c.rewardRatesMutex.Unlock()
	c.rewardRatesMapping[t] = rates
}

// GetLatestConfig returns latest system config
func (c *Cache) GetLatestConfig() models.Config {
	c.configMutex.RLock()
	defer c.configMutex.RUnlock()
	return c.config
}

// SetLatestConfig sets latest system config in cache
func (c *Cache) SetLatestConfig(config models.Config) {
	c.configMutex.Lock()
	defer c.configMutex.Unlock()
	c.config = config
}
