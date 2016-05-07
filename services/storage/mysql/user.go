package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/go-sql-driver/mysql"
)

// GetUserByID gets a user with id given
func (s Storage) GetUserByID(id int64) (models.User, error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE `id` = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrNotFound
		}

		return user, fmt.Errorf("query user by id error: %v", err)
	}

	return user, nil
}

// GetUserByEmail gets a user with email given
func (s Storage) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM users WHERE `email` = ?", email)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrNotFound
		}

		return user, fmt.Errorf("query user by email error: %v", err)
	}

	return user, nil
}

// CreateUser creates a new user
func (s Storage) CreateUser(u models.User) error {
	_, err := s.db.NamedExec("INSERT INTO users (`email`, `address`, `referer_id`) VALUES (:email, :address, :referer_id)", u)

	if err != nil {
		switch e := err.(type) {
		case *mysql.MySQLError:
			if e.Number == errcodeDuplicate {
				errcodeMapping := map[string]error{
					"key 'email'":   errors.ErrDuplicatedEmail,
					"key 'address'": errors.ErrDuplicatedAddress,
				}
				for k, v := range errcodeMapping {
					if strings.Contains(e.Message, k) {
						return v
					}
				}
			}
		default:
			return fmt.Errorf("create user error: %v", err)
		}
	}

	return nil
}

// UpdateUserStatus updates a user's status
func (s Storage) UpdateUserStatus(id int64, status string) error {
	_, err := s.db.Exec("UPDATE users SET `status` = ? WHERE `id` = ?", status, id)

	if err != nil {
		return fmt.Errorf("update user error: %v", err)
	}

	return nil
}

// GetRefereesSince gets user's referees since, id >= since
func (s Storage) GetRefereesSince(userID, id, limit int64) ([]models.User, error) {
	rawSQL := "SELECT * FROM users WHERE `referer_id` = ? AND `id` >= ? ORDER BY `id` ASC LIMIT ?"
	args := []interface{}{userID, id, limit}
	dest := []models.User{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetRefereesUntil gets user's referees since, id < until
func (s Storage) GetRefereesUntil(userID, id, limit int64) ([]models.User, error) {
	rawSQL := "SELECT * FROM users WHERE `referer_id` = ? AND `id` < ? ORDER BY `id` DESC LIMIT ?"
	args := []interface{}{userID, id, limit}
	dest := []models.User{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetWithdrawableUsers gets users who are able to withdraw
func (s Storage) GetWithdrawableUsers() ([]models.User, error) {
	rawSQL := "SELECT * FROM users WHERE `status` = ? AND `balance` > `min_withdrawal_amount`"
	args := []interface{}{models.UserStatusVerified}
	dest := []models.User{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}
