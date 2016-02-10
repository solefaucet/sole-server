package utils

import (
	"testing"

	"github.com/freeusd/solebtc/models"
)

func TestRandomReward(t *testing.T) {
	func() {
		defer func() {
			recover()
		}()
		RandomReward(nil)
	}()

	rates := []models.RewardRate{
		{
			Min:    0,
			Max:    10,
			Weight: 1,
		},
	}
	reward := RandomReward(rates)

	if reward >= 10 && reward < 0 {
		t.Errorf("reward should be [0, 10) but get %v", reward)
	}
}
