package mysql

import (
	"database/sql"
	"fmt"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetSessionByToken gets models.Session with token given
func (s Storage) GetSessionByToken(token string) (models.Session, *errors.Error) {
	session := models.Session{}
	err := s.db.Get(&session, "SELECT * FROM sessions WHERE token = ?", token)

	if err != nil {
		if err == sql.ErrNoRows {
			return session, errors.New(errors.ErrCodeNotFound)
		}

		return session, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get session unknown error: %v", err),
		}
	}

	return session, nil
}

// UpsertSession creates a new session
func (s Storage) UpsertSession(session models.Session) *errors.Error {
	_, err := s.db.NamedExec("INSERT INTO sessions (`user_id`, `token`, `type`) VALUES (:user_id, :token, :type) ON DUPLICATE KEY UPDATE `token` = :token", session)

	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Upsert session error: %v", err),
		}
	}

	return nil
}
