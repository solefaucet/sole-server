package mysql

import (
	"fmt"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// CreateWithdrawal creates a new withdrawal
func (s Storage) CreateWithdrawal(withdrawal models.Withdrawal) *errors.Error {
	tx := s.db.MustBegin()

	// create withdrawal with transaction
	if err := createWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create withdrawal commit transaction error: %v", err),
		}
	}

	return nil
}

func createWithdrawal(tx *sqlx.Tx, withdrawal models.Withdrawal) *errors.Error {
	if err := deductUserBalanceBy(tx, withdrawal.UserID, withdrawal.Amount); err != nil {
		return err
	}
	return insertWithdrawal(tx, withdrawal.UserID, withdrawal.BitcoinAddress, withdrawal.Amount)
}

func deductUserBalanceBy(tx *sqlx.Tx, userID int64, delta int64) *errors.Error {
	result, err := tx.Exec("UPDATE users SET `balance` = `balance` - ? WHERE `id` = ? AND `balance` >= ?", delta, userID, delta)
	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Deduct user balance error: %v", err),
		}
	}

	// make sure the user has sufficient balance
	if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return &errors.Error{
			ErrCode:             errors.ErrCodeInsufficientBalance,
			ErrStringForLogging: fmt.Sprintf("User balance is less than %v", delta),
		}
	}

	return nil
}

func insertWithdrawal(tx *sqlx.Tx, userID int64, bitcoinAddress string, amount int64) *errors.Error {
	rawSQL := "INSERT INTO withdrawals (`user_id`, `bitcoin_address`, `amount`) VALUES (?, ?, ?)"
	if _, err := tx.Exec(rawSQL, userID, bitcoinAddress, amount); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create withdrawal error error: %v", err),
		}
	}

	return nil
}

// GetWithdrawalsSince get user's withdrawal since, created_at >= since
func (s Storage) GetWithdrawalsSince(userID int64, since time.Time, limit int64) ([]models.Withdrawal, *errors.Error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `user_id` = ? AND `created_at` >= ? ORDER BY `id` ASC LIMIT ?"
	args := []interface{}{userID, since, limit}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetWithdrawalsUntil get user's withdrawal until, created_at < until
func (s Storage) GetWithdrawalsUntil(userID int64, until time.Time, limit int64) ([]models.Withdrawal, *errors.Error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `user_id` = ? AND `created_at` < ? ORDER BY `id` DESC LIMIT ?"
	args := []interface{}{userID, until, limit}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}
