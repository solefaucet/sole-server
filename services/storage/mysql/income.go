package mysql

import (
	"fmt"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// CreateRewardIncome creates a new reward type income
func (s Storage) CreateRewardIncome(userID, refererID, reward, rewardReferer int64, now time.Time) *errors.Error {
	income := models.Income{
		UserID:        userID,
		RefererID:     refererID,
		Type:          models.IncomeTypeReward,
		Income:        reward,
		RefererIncome: rewardReferer,
	}
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