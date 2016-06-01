package mysql

import (
	"database/sql"
	"fmt"

	"github.com/solefaucet/solebtc/models"
)

// GetLatestConfig get latest system config
func (s Storage) GetLatestConfig() (models.Config, error) {
	result := models.Config{}
	err := s.db.Get(&result, "SELECT * FROM configs ORDER BY `id` DESC LIMIT 1")

	if err != nil && err != sql.ErrNoRows {
		return result, fmt.Errorf("query latest config error: %v", err)
	}

	return result, nil
}
