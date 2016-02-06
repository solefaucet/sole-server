package mysql

import (
	"fmt"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetAuthToken gets models.AuthToken with auth_token given
func (s Storage) GetAuthToken(authToken string) (models.AuthToken, *errors.Error) {
	// TODO:
	return models.AuthToken{}, nil
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
	// TODO:
	return nil
}
