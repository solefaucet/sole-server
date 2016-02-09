package mysql

import (
	"fmt"
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetSortedTotalRewards get all total rewards order by time desc
func (s Storage) GetSortedTotalRewards() ([]models.TotalReward, *errors.Error) {
	trs := []models.TotalReward{}
	err := s.db.Select(&trs, "SELECT * FROM total_rewards ORDER BY created_at DESC")

	if err != nil {
		return trs, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get sorted total rewards unknown error: %v", err),
		}
	}

	return trs, nil
}

// IncrementTotalReward increments total reward by delta for now
func (s Storage) IncrementTotalReward(now time.Time, delta int64) *errors.Error {
	sql := "INSERT INTO total_rewards (`total`, `created_at`) VALUES (:delta, :created_at) ON DUPLICATE KEY UPDATE `total` = `total` + :delta"
	args := map[string]interface{}{
		"delta":      delta,
		"created_at": now,
	}
	_, err := s.db.NamedExec(sql, args)

	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Increment total reward unknown error: %v", err),
		}
	}

	return nil
}
