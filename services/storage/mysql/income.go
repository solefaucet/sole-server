package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

// GetRewardIncomes get user's reward incomes
func (s Storage) GetRewardIncomes(userID int64, limit, offset int64) ([]models.Income, error) {
	rawSQL := "SELECT * FROM incomes WHERE `user_id` = ? AND `type` = ? ORDER BY `id` DESC LIMIT ? OFFSET ?"
	args := []interface{}{userID, models.IncomeTypeReward, limit, offset}
	incomes := []models.Income{}
	err := s.selects(&incomes, rawSQL, args...)
	return incomes, err
}

// GetNumberOfRewardIncomes gets number of user's reward incomes
func (s Storage) GetNumberOfRewardIncomes(userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM incomes WHERE `user_id` = ? AND `type` = ?", userID, models.IncomeTypeReward).Scan(&count)
	return count, err
}

// CreateRewardIncome creates a new reward type income
func (s Storage) CreateRewardIncome(income models.Income, now time.Time) error {
	tx := s.db.MustBegin()

	if err := createRewardIncomeWithTx(tx, income, now); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create reward income commit transaction error: %v", err)
	}

	return nil
}

func createRewardIncomeWithTx(tx *sqlx.Tx, income models.Income, now time.Time) error {
	totalReward := income.Income

	// insert income into incomes table
	if err := insertIntoIncomesTableByRewardIncome(tx, income); err != nil {
		return err
	}

	// update user balance, total_income, referer_total_income
	if err := incrementUserBalanceByRewardIncome(tx, income.UserID, income.Income, income.RefererIncome, now); err != nil {
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

// update user balance, total_income, referer_total_income
func incrementUserBalanceByRewardIncome(tx *sqlx.Tx, userID int64, delta, refererDelta float64, now time.Time) error {
	rawSQL := "UPDATE users SET `balance` = `balance` + ?, `total_income` = `total_income` + ?, `referer_total_income` = `referer_total_income` + ?, `rewarded_at` = ? WHERE id = ?"
	args := []interface{}{delta, delta, refererDelta, now, userID}
	if result, err := tx.Exec(rawSQL, args...); err != nil {
		return fmt.Errorf("create reward income update user balance error: %v", err)
	} else if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("create reward income update user balance affected %v rows", rowAffected)
	}

	return nil
}

// update referer balance
func incrementRefererBalanceByRewardIncome(tx *sqlx.Tx, refererID int64, delta float64) (int64, error) {
	result, err := tx.NamedExec("UPDATE users SET `balance` = `balance` + :delta, `total_income_from_referees` = `total_income_from_referees` + :delta WHERE id = :id", map[string]interface{}{
		"id":    refererID,
		"delta": delta,
	})
	if err != nil {
		return 0, fmt.Errorf("update referer balance error: %v", err)
	}

	rowAffected, _ := result.RowsAffected()
	return rowAffected, nil
}

// insert reward income into incomes table
func insertIntoIncomesTableByRewardIncome(tx *sqlx.Tx, income models.Income) error {
	_, err := tx.NamedExec("INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (:user_id, :referer_id, :type, :income, :referer_income)", income)
	if err != nil {
		return fmt.Errorf("create reward income insert into incomes error: %v", err)
	}

	return nil
}

// increment total reward
func incrementTotalRewardByRewardIncome(tx *sqlx.Tx, totalReward float64, now time.Time) error {
	if _, err := tx.NamedExec(incrementTotalRewardSQL(now, totalReward)); err != nil {
		return fmt.Errorf("create reward income increment total reward error: %v", err)
	}

	return nil
}

// GetOfferwowEventByID finds offerwow event by event id
func (s Storage) GetOfferwowEventByID(eventID string) (models.OfferwowEvent, error) {
	event := models.OfferwowEvent{}
	err := s.db.Get(&event, "SELECT * FROM `offerwow` WHERE `event_id` = ?", eventID)

	if err != nil {
		if err == sql.ErrNoRows {
			return event, errors.ErrNotFound
		}

		return event, fmt.Errorf("query offerwow event by id error: %v", err)
	}

	return event, nil
}

// CreateOfferwowIncome creates a new reward type income
func (s Storage) CreateOfferwowIncome(income models.Income, eventID string) error {
	tx := s.db.MustBegin()

	if err := createOfferwowIncomeWithTx(tx, income, eventID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create reward income commit transaction error: %v", err)
	}

	return nil
}

func createOfferwowIncomeWithTx(tx *sqlx.Tx, income models.Income, eventID string) error {
	// update user balance, total_income, referer_total_income
	if err := incrementUserBalanceByOfferwowIncome(tx, income.UserID, income.Income, income.RefererIncome); err != nil {
		return err
	}

	// update referer balance
	if _, err := incrementRefererBalanceByRewardIncome(tx, income.RefererID, income.RefererIncome); err != nil {
		return err
	}

	// insert offerwow income into incomes table
	id, err := insertIntoIncomesTableByOfferwowIncome(tx, income)
	if err != nil {
		return err
	}

	// insert offerwow event
	offerwowEvent := models.OfferwowEvent{
		EventID:  eventID,
		IncomeID: id,
		Amount:   income.Income,
	}
	if _, err := tx.NamedExec("INSERT INTO `offerwow` (`event_id`, `income_id`, `amount`) VALUE (:event_id, :income_id, :amount)", offerwowEvent); err != nil {
		return err
	}

	return nil
}

// update user balance, total_income, referer_total_income
func incrementUserBalanceByOfferwowIncome(tx *sqlx.Tx, userID int64, delta, refererDelta float64) error {
	rawSQL := "UPDATE users SET `balance` = `balance` + ?, `total_income` = `total_income` + ?, `referer_total_income` = `referer_total_income` + ? WHERE id = ?"
	args := []interface{}{delta, delta, refererDelta, userID}
	if result, err := tx.Exec(rawSQL, args...); err != nil {
		return fmt.Errorf("create offerwow income update user balance error: %v", err)
	} else if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("create offerwow income update user balance affected %v rows", rowAffected)
	}

	return nil
}

// insert offerwow income into incomes table
func insertIntoIncomesTableByOfferwowIncome(tx *sqlx.Tx, income models.Income) (int64, error) {
	result, err := tx.NamedExec("INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (:user_id, :referer_id, :type, :income, :referer_income)", income)
	if err != nil {
		return 0, fmt.Errorf("create offerwow income insert into incomes error: %v", err)
	}

	return result.LastInsertId()
}
