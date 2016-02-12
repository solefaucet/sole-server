package memory

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/freeusd/solebtc/models"
)

func mockGetBitcoinPriceFunc(p int64, err error) getBitcoinPriceFunc {
	return func() (int64, error) {
		return p, err
	}
}

func TestMemory(t *testing.T) {
	c := New(mockGetBitcoinPriceFunc(8, nil), &bytes.Buffer{}, time.Second)

	p := c.GetBitcoinPrice()
	if p != 8 {
		t.Errorf("price should be 8.8 but get %v", p)
	}

	c.getBTCPrice = mockGetBitcoinPriceFunc(0, nil)
	c.setBitcoinPrice(false)

	c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New(""))
	c.setBitcoinPrice(false)

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, nil)
		c.setBitcoinPrice(true)
	})

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New(""))
		c.setBitcoinPrice(true)
	})

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New("error"))
		c.backgroundJob(true, time.Second)
	})

	now := time.Now()

	c.IncrementTotalReward(now, 1)
	if v := c.GetLatestTotalReward(); v.Total != 1 {
		t.Errorf("total reward should be 1 but get %v", v)
	}

	c.IncrementTotalReward(now, 1)
	if v := c.GetLatestTotalReward(); v.Total != 2 {
		t.Errorf("total reward should be 2 but get %v", v)
	}

	c.SetRewardRates(models.RewardRateTypeLess, []models.RewardRate{models.RewardRate{}})
	rates := c.GetRewardRatesByType(models.RewardRateTypeLess)
	if len(rates) != 1 {
		t.Errorf("expected length of rates should be 1 but get %v", len(rates))
	}

	c.SetLatestConfig(models.Config{TotalRewardThreshold: 1000})
	config := c.GetLatestConfig()
	if config.TotalRewardThreshold != 1000 {
		t.Errorf("expected total reward threshold should be 1000 but get %v", config.TotalRewardThreshold)
	}
}

func funcWithRecover(f func()) {
	defer func() {
		recover()
	}()
	f()
}
