package memory

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/freeusd/solebtc/models"
)

func TestMemory(t *testing.T) {
	c := New(20)

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

	c.UpdateBitcoinPrice(1)
	bitcoinPrice := c.GetLatestConfig().BitcoinPrice
	if bitcoinPrice != 1 {
		t.Errorf("expected bitcoin price should be 1 but get %v", bitcoinPrice)
	}

	c.InsertIncome(map[string]interface{}{"k": 1})
	c.InsertIncome(map[string]interface{}{"k": 2})
	expectedLatestIncomes := `[{"k":2},{"k":1}]`
	actualLatestIncomes, _ := json.Marshal(c.GetLatestIncomes())
	if actual := string(actualLatestIncomes); actual != expectedLatestIncomes {
		t.Errorf("expected bytes %v but get %v", expectedLatestIncomes, actual)
	}
}

func BenchmarkInsertIncome(b *testing.B) {
	c := New(20)
	object := map[string]interface{}{
		"key1": 123.4123,
		"key2": "sdflkjei",
		"key3": time.Now(),
		"key4": "bitcoin_address_very_long_string",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.InsertIncome(object)
	}
}
