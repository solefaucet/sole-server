package mysql

import (
	"fmt"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// CreateWithdrawl creates a new withdrawl
func (s Storage) CreateWithdrawl(withdrawl models.Withdrawl) *errors.Error {
	tx := s.db.MustBegin()

	// create withdrawl with transaction
	if err := createWithdrawl(tx, withdrawl); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create withdrawl commit transaction error: %v", err),
		}
	}

	return nil
}

func createWithdrawl(tx *sqlx.Tx, withdrawl models.Withdrawl) *errors.Error {
	if err := deductUserBalanceBy(tx, withdrawl.UserID, withdrawl.Amount); err != nil {
		return err
	}
	return insertWithdrawl(tx, withdrawl.UserID, withdrawl.BitcoinAddress, withdrawl.Amount)
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

func insertWithdrawl(tx *sqlx.Tx, userID int64, bitcoinAddress string, amount int64) *errors.Error {
	rawSQL := "INSERT INTO withdrawls (`user_id`, `bitcoin_address`, `amount`) VALUES (?, ?, ?)"
	if _, err := tx.Exec(rawSQL, userID, bitcoinAddress, amount); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create withdrawl error error: %v", err),
		}
	}

	return nil
}
