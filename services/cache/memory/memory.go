package memory

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/freeusd/solebtc/services/cache"
)

// Cache implements cache.Cache interface with memory
type Cache struct {
	getBTCPrice    getBitcoinPriceFunc
	cachedBTCPrice float64

	logWriter io.Writer
}

var _ cache.Cache = Cache{}

type getBitcoinPriceFunc func() (float64, error)

// New creates a new in-memory cache
func New(getBTCPrice getBitcoinPriceFunc, logWriter io.Writer, interval time.Duration) Cache {
	c := Cache{
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
func (c Cache) GetBitcoinPrice() float64 {
	return c.cachedBTCPrice
}

var epsilon = math.Nextafter(1, 2) - 1

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
	p, err := c.getBTCPrice()
	if err != nil {
		if panicOnError {
			panic(err)
		}
		fmt.Fprintf(c.logWriter, "get bitcoin price error: %v", err)
		return
	}

	if -epsilon < p && p < epsilon {
		errorString := fmt.Sprintf("bitcoin price %v should not be 0", p)
		if panicOnError {
			panic(errorString)
		}
		fmt.Fprintf(c.logWriter, errorString)
		return
	}

	c.cachedBTCPrice = p
}
