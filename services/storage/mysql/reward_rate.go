package mysql

import (
	"fmt"

	"github.com/solefaucet/solebtc/models"
)

// GetRewardRatesByType get all reward rates by type
func (s Storage) GetRewardRatesByType(rewardRateType string) ([]models.RewardRate, error) {
	rrs := []models.RewardRate{}
	err := s.db.Select(&rrs, "SELECT * FROM reward_rates WHERE `type` = ?", rewardRateType)

	if err != nil {
		return nil, fmt.Errorf("query reward rates error: %v", err)
	}

	return rrs, nil
}
