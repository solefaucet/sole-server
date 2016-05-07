package mysql

import (
	"database/sql"
	"fmt"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/go-sql-driver/mysql"
)

// GetAuthToken gets models.AuthToken with auth_token given
func (s Storage) GetAuthToken(authTokenString string) (models.AuthToken, *errors.Error) {
	authToken := models.AuthToken{}
	err := s.db.Get(&authToken, "SELECT * FROM auth_tokens WHERE auth_token = ?", authTokenString)

	if err != nil {
		if err == sql.ErrNoRows {
			return authToken, errors.New(errors.ErrCodeNotFound)
		}

		return authToken, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get auth token unknown error: %v", err),
		}
	}

	return authToken, nil
}

// CreateAuthToken creates a new auth token
func (s Storage) CreateAuthToken(authToken models.AuthToken) *errors.Error {
	_, err := s.db.NamedExec("INSERT INTO auth_tokens (`user_id`, `auth_token`) VALUES (:user_id, :auth_token)", authToken)

	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == errcodeDuplicate {
				return errors.New(errors.ErrCodeDuplicateAuthToken)
			}
		}

		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create auth token unknown error: %v", err),
		}
	}

	return nil
}

// DeleteAuthToken deletes auth_token from storage
func (s Storage) DeleteAuthToken(authToken string) *errors.Error {
	_, err := s.db.Exec("DELETE FROM auth_tokens WHERE auth_token = ?", authToken)

	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Delete auth token unknown error: %v", err),
		}
	}

	return nil
}
