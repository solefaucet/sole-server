package mysql

import (
	"database/sql"
	"fmt"

	"github.com/solefaucet/sole-server/models"
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
