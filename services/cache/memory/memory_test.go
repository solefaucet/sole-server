package memory

import (
	"testing"
	"time"

	"github.com/freeusd/solebtc/models"
)

func TestMemory(t *testing.T) {
	c := New()

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
