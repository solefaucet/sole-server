package mysql

import (
	"fmt"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// GetRewardIncomesSince get user's reward income records since, created_at >= since
// pagination design, previous
// https://developers.facebook.com/blog/post/478/
func (s Storage) GetRewardIncomesSince(userID int64, since time.Time, limit int64) ([]models.Income, *errors.Error) {
	rawSQL := "SELECT * FROM incomes WHERE `user_id` = ? AND `created_at` >= ? AND `type` = ? ORDER BY `id` ASC LIMIT ?"
	args := []interface{}{userID, since, models.IncomeTypeReward, limit}
	return s.getIncomes(rawSQL, args...)
}

// GetRewardIncomesUntil get user's reward income records until, created_at < until
// pagination design, next
func (s Storage) GetRewardIncomesUntil(userID int64, until time.Time, limit int64) ([]models.Income, *errors.Error) {
	rawSQL := "SELECT * FROM incomes WHERE `user_id` = ? AND `created_at` < ? AND `type` = ? ORDER BY `id` DESC LIMIT ?"
	args := []interface{}{userID, until, models.IncomeTypeReward, limit}
	return s.getIncomes(rawSQL, args...)
}

func (s Storage) getIncomes(rawSQL string, args ...interface{}) ([]models.Income, *errors.Error) {
	incomes := []models.Income{}
	if err := s.db.Select(&incomes, rawSQL, args...); err != nil {
		return nil, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Get reward incomes unknown error: %v", err),
		}
	}

	return incomes, nil
}

// CreateRewardIncome creates a new reward type income
func (s Storage) CreateRewardIncome(income models.Income, now time.Time) *errors.Error {
	tx := s.db.MustBegin()

	if err := createRewardIncomeWithTx(tx, income, now); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income transaction commit unknown error: %v", err),
		}
	}

	return nil
}

func createRewardIncomeWithTx(tx *sqlx.Tx, income models.Income, now time.Time) *errors.Error {
	totalReward := income.Income

	// insert income into incomes table
	if err := insertIntoIncomesTableByRewardIncome(tx, income); err != nil {
		return err
	}

	// update user balance
	if err := incrementUserBalanceByRewardIncome(tx, income.UserID, income.Income, now); err != nil {
		return err
	}

	// update referer balance
	if rowAffected, err := incrementRefererBalanceByRewardIncome(tx, income.RefererID, income.RefererIncome); err != nil {
		return err
	} else if rowAffected == 1 {
		totalReward += income.RefererIncome
	}

	// update total reward
	if err := incrementTotalRewardByRewardIncome(tx, totalReward, now); err != nil {
		return err
	}

	return nil
}

// update user balance
func incrementUserBalanceByRewardIncome(tx *sqlx.Tx, userID int64, delta int64, now time.Time) *errors.Error {
	if result, err := tx.Exec("UPDATE users SET `balance` = `balance` + ?, `rewarded_at` = ? WHERE id = ?", delta, now, userID); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income update user balance unknown error: %v", err),
		}
	} else if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income update user balance affected %v rows", rowAffected),
		}
	}

	return nil
}

// update referer balance
func incrementRefererBalanceByRewardIncome(tx *sqlx.Tx, refererID int64, delta int64) (int64, *errors.Error) {
	result, err := tx.Exec("UPDATE users SET `balance` = `balance` + ? WHERE id = ?", delta, refererID)
	if err != nil {
		return 0, &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income update referer balance unknown error: %v", err),
		}
	}

	rowAffected, _ := result.RowsAffected()
	return rowAffected, nil
}

// insert into incomes table
func insertIntoIncomesTableByRewardIncome(tx *sqlx.Tx, income models.Income) *errors.Error {
	_, err := tx.NamedExec("INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (:user_id, :referer_id, :type, :income, :referer_income)", income)
	if err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income insert into incomes unknown error: %v", err),
		}
	}

	return nil
}

// increment total reward
func incrementTotalRewardByRewardIncome(tx *sqlx.Tx, totalReward int64, now time.Time) *errors.Error {
	if _, err := tx.NamedExec(incrementTotalRewardSQL(now, totalReward)); err != nil {
		return &errors.Error{
			ErrCode:             errors.ErrCodeUnknown,
			ErrStringForLogging: fmt.Sprintf("Create reward income increment total reward unknown error: %v", err),
		}
	}

	return nil
}
