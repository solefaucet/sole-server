package mysql

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

// GetAuthToken gets models.AuthToken with auth_token given
func (s Storage) GetAuthToken(authTokenString string) (models.AuthToken, error) {
	authToken := models.AuthToken{}
	err := s.db.Get(&authToken, "SELECT * FROM auth_tokens WHERE auth_token = ?", authTokenString)

	if err != nil {
		if err == sql.ErrNoRows {
			return authToken, errors.ErrNotFound
		}

		return authToken, fmt.Errorf("query auth token error: %v", err)
	}

	return authToken, nil
}

// CreateAuthToken creates a new auth token
func (s Storage) CreateAuthToken(authToken models.AuthToken) error {
	_, err := s.db.NamedExec("INSERT INTO auth_tokens (`user_id`, `auth_token`) VALUES (:user_id, :auth_token)", authToken)

	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == errcodeDuplicate {
				return errors.ErrDuplicatedAuthToken
			}
		}

		return fmt.Errorf("create auth token error: %v", err)
	}

	return nil
}

// DeleteAuthToken deletes auth_token from storage
func (s Storage) DeleteAuthToken(authToken string) error {
	_, err := s.db.Exec("DELETE FROM auth_tokens WHERE auth_token = ?", authToken)

	if err != nil {
		return fmt.Errorf("delete auth token error: %v", err)
	}

	return nil
}
