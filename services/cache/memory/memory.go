package memory

import (
	"container/ring"
	"sync"
	"time"

	"github.com/solefaucet/sole-server/models"
)

// Cache implements cache.Cache interface with memory
type Cache struct {
	totalReward      models.TotalReward
	totalRewardMutex sync.RWMutex

	rewardRatesMapping map[string][]models.RewardRate
	rewardRatesMutex   sync.RWMutex

	config      models.Config
	configMutex sync.RWMutex

	incomesRing   *ring.Ring
	latestIncomes []interface{}
	incomesMutex  sync.RWMutex
}

// var _ cache.Cache = &Cache{}

// New creates a new in-memory cache
// number of cached incomes
func New(numCachedIncomes int) *Cache {
	return &Cache{
		rewardRatesMapping: make(map[string][]models.RewardRate),
		incomesRing:        ring.New(numCachedIncomes),
	}
}

// GetLatestTotalReward returns total reward of today
func (c *Cache) GetLatestTotalReward() models.TotalReward {
	c.totalRewardMutex.RLock()
	defer c.totalRewardMutex.RUnlock()
	return c.totalReward
}

// IncrementTotalReward increment total reward today by delta if day matches
func (c *Cache) IncrementTotalReward(t time.Time, delta float64) {
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

// InsertIncome inserts a new income
func (c *Cache) InsertIncome(income interface{}) {
	c.incomesMutex.Lock()
	defer c.incomesMutex.Unlock()
	// put it in head
	c.incomesRing = c.incomesRing.Prev()
	c.incomesRing.Value = income

	// re-calculate bytes
	c.latestIncomes = []interface{}{}
	c.incomesRing.Do(func(i interface{}) {
		if i != nil {
			c.latestIncomes = append(c.latestIncomes, i)
		}
	})
}

// GetLatestIncomes returns incomes
func (c *Cache) GetLatestIncomes() []interface{} {
	c.incomesMutex.RLock()
	defer c.incomesMutex.RUnlock()
	return c.latestIncomes
}
