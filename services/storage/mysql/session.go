package mysql

import (
	"database/sql"
	"fmt"

	"github.com/solefaucet/solebtc/errors"
	"github.com/solefaucet/solebtc/models"
)

// GetSessionByToken gets models.Session with token given
func (s Storage) GetSessionByToken(token string) (models.Session, error) {
	session := models.Session{}
	err := s.db.Get(&session, "SELECT * FROM sessions WHERE token = ?", token)

	if err != nil {
		if err == sql.ErrNoRows {
			return session, errors.ErrNotFound
		}

		return session, fmt.Errorf("query session error: %v", err)
	}

	return session, nil
}

// UpsertSession creates a new session
func (s Storage) UpsertSession(session models.Session) error {
	_, err := s.db.NamedExec("INSERT INTO sessions (`user_id`, `token`, `type`) VALUES (:user_id, :token, :type) ON DUPLICATE KEY UPDATE `token` = :token", session)

	if err != nil {
		return fmt.Errorf("upsert session error: %v", err)
	}

	return nil
}
