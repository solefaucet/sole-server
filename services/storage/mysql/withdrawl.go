package mysql

import (
	"fmt"
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/jmoiron/sqlx"
)

// CreateWithdrawal creates a new withdrawal
func (s Storage) CreateWithdrawal(withdrawal models.Withdrawal) error {
	tx := s.db.MustBegin()

	// create withdrawal with transaction
	if err := createWithdrawal(tx, withdrawal); err != nil {
		tx.Rollback()
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create withdrawal commit transaction error: %v", err)
	}

	return nil
}

func createWithdrawal(tx *sqlx.Tx, withdrawal models.Withdrawal) error {
	if err := deductUserBalanceBy(tx, withdrawal.UserID, withdrawal.Amount); err != nil {
		return err
	}
	return insertWithdrawal(tx, withdrawal.UserID, withdrawal.Address, withdrawal.Amount)
}

func deductUserBalanceBy(tx *sqlx.Tx, userID int64, delta float64) error {
	result, err := tx.Exec("UPDATE users SET `balance` = `balance` - ? WHERE `id` = ? AND `balance` >= ?", delta, userID, delta)
	if err != nil {
		return fmt.Errorf("deduct user balance error: %v", err)
	}

	// make sure the user has sufficient balance
	if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return errors.ErrInsufficientBalance
	}

	return nil
}

func insertWithdrawal(tx *sqlx.Tx, userID int64, address string, amount float64) error {
	rawSQL := "INSERT INTO withdrawals (`user_id`, `address`, `amount`) VALUES (?, ?, ?)"
	if _, err := tx.Exec(rawSQL, userID, address, amount); err != nil {
		return fmt.Errorf("create withdrawal error: %v", err)
	}

	return nil
}

// GetWithdrawalsSince get user's withdrawal since, created_at >= since
func (s Storage) GetWithdrawalsSince(userID int64, since time.Time, limit int64) ([]models.Withdrawal, error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `user_id` = ? AND `created_at` >= ? ORDER BY `id` ASC LIMIT ?"
	args := []interface{}{userID, since, limit}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetWithdrawalsUntil get user's withdrawal until, created_at < until
func (s Storage) GetWithdrawalsUntil(userID int64, until time.Time, limit int64) ([]models.Withdrawal, error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `user_id` = ? AND `created_at` < ? ORDER BY `id` DESC LIMIT ?"
	args := []interface{}{userID, until, limit}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetUnprocessedWithdrawals get all unprocessed withdrawals
func (s Storage) GetUnprocessedWithdrawals() ([]models.Withdrawal, error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `status` != ? ORDER BY `id` asc"
	args := []interface{}{models.WithdrawalStatusProcessed}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// UpdateWithdrawalStatusToProcessing update withdrawal status to processing if status = pending
func (s Storage) UpdateWithdrawalStatusToProcessing(id int64) error {
	rawSQL := "UPDATE `withdrawals` SET `status` = ? WHERE `id` = ? AND `status` = ?"
	args := []interface{}{models.WithdrawalStatusProcessing, id, models.WithdrawalStatusPending}
	result, err := s.db.Exec(rawSQL, args...)
	if err != nil {
		return err
	}

	if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("expected 1 but %v rows affected", rowAffected)
	}

	return nil
}

// UpdateWithdrawalStatusToProcessed update withdrawal status to processed if status = processing
func (s Storage) UpdateWithdrawalStatusToProcessed(id int64, transactionID string) error {
	rawSQL := "UPDATE `withdrawals` SET `status` = ?, `transaction_id` = ? WHERE `id` = ? AND `status` = ?"
	args := []interface{}{models.WithdrawalStatusProcessed, transactionID, id, models.WithdrawalStatusProcessing}
	result, err := s.db.Exec(rawSQL, args...)
	if err != nil {
		return err
	}

	if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("expected 1 but %v rows affected", rowAffected)
	}

	return nil
}
