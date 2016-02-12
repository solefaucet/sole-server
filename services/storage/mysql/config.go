package mysql

import (
	"database/sql"
	"fmt"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetLatestConfig get latest system config
func (s Storage) GetLatestConfig() (models.Config, *errors.Error) {
	result := models.Config{}
	err := s.db.Get(&result, "SELECT * FROM configs ORDER BY `id` DESC LIMIT 1")

	if err != nil && err != sql.ErrNoRows {
		return result, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get latest config unknown error: %v", err),
		}
	}

	return result, nil
}
