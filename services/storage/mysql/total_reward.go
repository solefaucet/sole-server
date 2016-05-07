package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/freeusd/solebtc/models"
)

// GetLatestTotalReward get all total rewards order by time desc
func (s Storage) GetLatestTotalReward() (models.TotalReward, error) {
	result := models.TotalReward{}
	err := s.db.Get(&result, "SELECT * FROM total_rewards FORCE INDEX (`created_at`) ORDER BY `created_at` DESC LIMIT 1")

	if err != nil && err != sql.ErrNoRows {
		return result, fmt.Errorf("query sorted total rewards error: %v", err)
	}

	return result, nil
}

// IncrementTotalReward increments total reward by delta for now
func (s Storage) IncrementTotalReward(now time.Time, delta float64) error {
	_, err := s.db.NamedExec(incrementTotalRewardSQL(now, delta))

	if err != nil {
		return fmt.Errorf("increment total reward error: %v", err)
	}

	return nil
}

func incrementTotalRewardSQL(now time.Time, delta float64) (string, map[string]interface{}) {
	sql := "INSERT INTO total_rewards (`total`, `created_at`) VALUES (:delta, :created_at) ON DUPLICATE KEY UPDATE `total` = `total` + :delta"
	args := map[string]interface{}{
		"delta":      delta,
		"created_at": now,
	}
	return sql, args
}
