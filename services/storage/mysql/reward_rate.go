package mysql

import (
	"fmt"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetRewardRatesByType get all reward rates by type
func (s Storage) GetRewardRatesByType(rewardRateType string) ([]models.RewardRate, *errors.Error) {
	rrs := []models.RewardRate{}
	err := s.db.Select(&rrs, "SELECT * FROM reward_rates WHERE `type` = ?", rewardRateType)

	if err != nil {
		return rrs, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get reward rates unknown error: %v", err),
		}
	}

	return rrs, nil
}
