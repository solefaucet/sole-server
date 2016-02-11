package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/freeusd/solebtc/models"
)

// RandomReward generates a random reward with rates given
func RandomReward(rates []models.RewardRate) int64 {
	var sum int64
	for i := range rates {
		sum += rates[i].Weight
	}
	if sum < 1 {
		panic("sum of reward rates weight should be greater than 0")
	}

	i := 0
	for r := randInt64(0, sum); i < len(rates); i++ {
		r -= rates[i].Weight
		if r < 0 {
			break
		}
	}

	rate := rates[i]
	return randInt64(rate.Min, rate.Max)
}

func randInt64(min, max int64) int64 {
	// panic if rand.Int returns error, fail fast here
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return min + n.Int64()
}
