package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetUserByID gets a user with id given
func (s Storage) GetUserByID(id int64) (models.User, *errors.Error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE `id` = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New(errors.ErrCodeNotFound)
		}

		return user, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get user unknown error: %v", err),
		}
	}

	return user, nil
}

// GetUserByEmail gets a user with email given
func (s Storage) GetUserByEmail(email string) (models.User, *errors.Error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE `email` = ?", email)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New(errors.ErrCodeNotFound)
		}

		return user, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get user unknown error: %v", err),
		}
	}

	return user, nil
}

// CreateUser creates a new user
func (s Storage) CreateUser(u models.User) *errors.Error {
	_, err := s.db.NamedExec("INSERT INTO users (`email`, `bitcoin_address`, `referer_id`) VALUES (:email, :bitcoin_address, :referer_id)", u)

	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == errcodeDuplicate {
				syserr := errors.New(errors.ErrCodeUnknown)
				errcodeMapping := map[string]errors.Code{
					"key 'email'":           errors.ErrCodeDuplicateEmail,
					"key 'bitcoin_address'": errors.ErrCodeDuplicateBitcoinAddress,
				}
				for k, v := range errcodeMapping {
					if strings.Contains(e.Message, k) {
						syserr.ErrCode = v
					}
				}
				return syserr
			}
		}

		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create user unknown error: %v", err),
		}
	}

	return nil
}

// UpdateUser updates a user's info
func (s Storage) UpdateUser(user models.User) *errors.Error {
	_, err := s.db.NamedExec("UPDATE users SET `status` = :status WHERE `id` = :id", user)

	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Update user unknown error: %v", err),
		}
	}

	return nil
}

// IncrementUserBalance increments user's balance by delta and update rewarded_at if needed
func (s Storage) IncrementUserBalance(userID int64, delta int64, t *time.Time) *errors.Error {
	var err error
	if t != nil {
		_, err = s.db.Exec("UPDATE users SET `balance` = `balance` + ?, `rewarded_at` = ? WHERE id = ?", delta, *t, userID)
	} else {
		_, err = s.db.Exec("UPDATE users SET `balance` = `balance` + ? WHERE id = ?", delta, userID)
	}

	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Increment user balance unknown error: %v", err),
		}
	}

	return nil
}
