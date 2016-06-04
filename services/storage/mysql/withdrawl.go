package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
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

// GetWithdrawals get user's withdrawal
func (s Storage) GetWithdrawals(userID int64, limit, offset int64) ([]models.Withdrawal, error) {
	rawSQL := "SELECT * FROM `withdrawals` WHERE `user_id` = ? ORDER BY `id` DESC LIMIT ? OFFSET ?"
	args := []interface{}{userID, limit, offset}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// GetNumberOfWithdrawals gets number of user's withdrawals
func (s Storage) GetNumberOfWithdrawals(userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `withdrawals` WHERE `user_id` = ?", userID).Scan(&count)
	return count, err
}

// GetPendingWithdrawals get all unprocessed withdrawals
func (s Storage) GetPendingWithdrawals() ([]models.Withdrawal, error) {
	rawSQL := "SELECT * FROM withdrawals WHERE `status` = ? ORDER BY `id` asc"
	args := []interface{}{models.WithdrawalStatusPending}
	dest := []models.Withdrawal{}
	err := s.selects(&dest, rawSQL, args...)
	return dest, err
}

// UpdateWithdrawalStatusToProcessing update withdrawal status to processing if status = pending
func (s Storage) UpdateWithdrawalStatusToProcessing(ids []int64) error {
	rawSQL, args, err := sqlx.In(
		"UPDATE `withdrawals` SET `status` = ? WHERE `id` IN (?) AND `status` = ?",
		models.WithdrawalStatusProcessing,
		ids,
		models.WithdrawalStatusPending,
	)
	if err != nil {
		return fmt.Errorf("update withdrawal status to processing build sql with in: %v", err)
	}

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
func (s Storage) UpdateWithdrawalStatusToProcessed(ids []int64, transactionID string) error {
	rawSQL, args, err := sqlx.In(
		"UPDATE `withdrawals` SET `status` = ?, `transaction_id` = ? WHERE `id` IN (?) AND `status` = ?",
		models.WithdrawalStatusProcessed,
		transactionID,
		ids,
		models.WithdrawalStatusProcessing,
	)
	if err != nil {
		return fmt.Errorf("update withdrawal status to processed build sql with in: %v", err)
	}

	result, err := s.db.Exec(rawSQL, args...)
	if err != nil {
		return err
	}

	if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("expected 1 but %v rows affected", rowAffected)
	}

	return nil
}
