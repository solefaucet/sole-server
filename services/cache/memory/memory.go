package memory

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/freeusd/solebtc/services/cache"
)

// Cache implements cache.Cache interface with memory
type Cache struct {
	getBTCPrice         getBitcoinPriceFunc
	cachedBTCPrice      int64
	cachedBTCPriceMutex sync.RWMutex

	lastDay          int
	totalReward      int64
	totalRewardMutex sync.RWMutex

	logWriter io.Writer
}

var _ cache.Cache = &Cache{}

type getBitcoinPriceFunc func() (int64, error)

// New creates a new in-memory cache
func New(getBTCPrice getBitcoinPriceFunc, logWriter io.Writer, interval time.Duration) *Cache {
	c := &Cache{
		getBTCPrice: getBTCPrice,
		logWriter:   logWriter,
	}

	// get init value, on error it should panic
	c.setBitcoinPrice(true)

	// get new bitcoin price in the background
	go c.backgroundJob(false, interval)

	return c
}

// GetBitcoinPrice returns the cached bitcoin price
func (c *Cache) GetBitcoinPrice() int64 {
	c.cachedBTCPriceMutex.RLock()
	defer c.cachedBTCPriceMutex.RUnlock()
	return c.cachedBTCPrice
}

// GetTotalRewardToday returns total reward of today
func (c *Cache) GetTotalRewardToday() int64 {
	c.totalRewardMutex.RLock()
	defer c.totalRewardMutex.RUnlock()
	return c.totalReward
}

// IncrementTotalReward increment total reward today by delta if day matches
func (c *Cache) IncrementTotalReward(t time.Time, delta int64) {
	c.totalRewardMutex.Lock()
	defer c.totalRewardMutex.Unlock()

	if c.lastDay == t.YearDay() {
		c.totalReward += delta
	} else {
		c.totalReward = delta
		c.lastDay = t.YearDay()
	}
}

func (c *Cache) backgroundJob(panicOnError bool, interval time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(c.logWriter, "memory cache background job panic: %v", err)
		}
	}()

	for {
		select {
		case <-time.After(interval):
			c.setBitcoinPrice(panicOnError)
		}
	}
}

func (c *Cache) setBitcoinPrice(panicOnError bool) {
	c.cachedBTCPriceMutex.Lock()
	defer c.cachedBTCPriceMutex.Unlock()

	p, err := c.getBTCPrice()
	if err != nil {
		if panicOnError {
			panic(err)
		}
		fmt.Fprintf(c.logWriter, "get bitcoin price error: %v", err)
		return
	}

	if p == 0 {
		errorString := fmt.Sprintf("bitcoin price %v should not be 0", p)
		if panicOnError {
			panic(errorString)
		}
		fmt.Fprintf(c.logWriter, errorString)
		return
	}

	c.cachedBTCPrice = p
}
