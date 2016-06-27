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
	if _, err := addIncome(tx, income); err != nil {
		return err
	}

	// update user balance, total_income, referer_total_income
	if err := incrementUserBalance(tx, income.UserID, income.Income, income.RefererIncome); err != nil {
		return err
	}

	// update user rewarded_at
	if _, err := tx.Exec("UPDATE users SET `rewarded_at` = ? WHERE `id` = ?", now, income.UserID); err != nil {
		return err
	}

	// update referer balance
	if rowAffected, err := incrementRefererBalance(tx, income.RefererID, income.RefererIncome); err != nil {
		return err
	} else if rowAffected == 1 {
		totalReward += income.RefererIncome
	}

	// update total reward
	if err := incrementTotalReward(tx, totalReward, now); err != nil {
		return err
	}

	return nil
}

// insert reward income into incomes table
func addIncome(tx *sqlx.Tx, income models.Income) (sql.Result, error) {
	result, err := tx.NamedExec("INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (:user_id, :referer_id, :type, :income, :referer_income)", income)
	if err != nil {
		return nil, fmt.Errorf("add income error: %v", err)
	}

	return result, nil
}

// increment user balance, total_income, referer_total_income
func incrementUserBalance(tx *sqlx.Tx, userID int64, delta, refererDelta float64) error {
	rawSQL := "UPDATE users SET `balance` = `balance` + ?, `total_income` = `total_income` + ?, `referer_total_income` = `referer_total_income` + ? WHERE id = ?"
	args := []interface{}{delta, delta, refererDelta, userID}
	if result, err := tx.Exec(rawSQL, args...); err != nil {
		return fmt.Errorf("increment user balance error: %v", err)
	} else if rowAffected, _ := result.RowsAffected(); rowAffected != 1 {
		return fmt.Errorf("increment user balance affected %v rows", rowAffected)
	}

	return nil
}

// increment referer balance
func incrementRefererBalance(tx *sqlx.Tx, refererID int64, delta float64) (int64, error) {
	result, err := tx.NamedExec("UPDATE users SET `balance` = `balance` + :delta, `total_income_from_referees` = `total_income_from_referees` + :delta WHERE id = :id", map[string]interface{}{
		"id":    refererID,
		"delta": delta,
	})
	if err != nil {
		return 0, fmt.Errorf("increment referer balance error: %v", err)
	}

	rowAffected, _ := result.RowsAffected()
	return rowAffected, nil
}

// increment total reward
func incrementTotalReward(tx *sqlx.Tx, totalReward float64, now time.Time) error {
	sql := "INSERT INTO total_rewards (`total`, `created_at`) VALUES (:delta, :created_at) ON DUPLICATE KEY UPDATE `total` = `total` + :delta"
	args := map[string]interface{}{
		"delta":      totalReward,
		"created_at": now,
	}

	if _, err := tx.NamedExec(sql, args); err != nil {
		return fmt.Errorf("increment total reward error: %v", err)
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

// CreateOfferwowIncome creates a new offerwow type income
func (s Storage) CreateOfferwowIncome(income models.Income, eventID string) error {
	tx := s.db.MustBegin()

	if err := createOfferwowIncomeWithTx(tx, income, eventID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create offerwow income commit transaction error: %v", err)
	}

	return nil
}

func createOfferwowIncomeWithTx(tx *sqlx.Tx, income models.Income, eventID string) error {
	// insert offerwow income into incomes table
	result, err := addIncome(tx, income)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// update user balance, total_income, referer_total_income
	if err := incrementUserBalance(tx, income.UserID, income.Income, income.RefererIncome); err != nil {
		return err
	}

	// update referer balance
	if _, err := incrementRefererBalance(tx, income.RefererID, income.RefererIncome); err != nil {
		return err
	}

	// insert offerwow event
	offerwowEvent := models.OfferwowEvent{
		EventID:  eventID,
		IncomeID: id,
		Amount:   income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `offerwow` (`event_id`, `income_id`, `amount`) VALUE (:event_id, :income_id, :amount)", offerwowEvent)
	return err
}

// GetSuperrewardsOfferByID finds superrewards offer by transaction Id and user id
func (s Storage) GetSuperrewardsOfferByID(transactionID string, userID int64) (models.SuperrewardsOffer, error) {
	offer := models.SuperrewardsOffer{}
	err := s.db.Get(&offer, "SELECT * FROM `superrewards` WHERE `transaction_id` = ? AND `user_id` = ?", transactionID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return offer, errors.ErrNotFound
		}

		return offer, fmt.Errorf("query superrewards offer by id error: %v", err)
	}

	return offer, nil
}

// CreateSuperrewardsIncome creates a new superrewards type income
func (s Storage) CreateSuperrewardsIncome(income models.Income, transactionID, offerID string) error {
	tx := s.db.MustBegin()

	if err := createSuperrewardsIncomeWithTx(tx, income, transactionID, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create superrewards income commit transaction error: %v", err)
	}

	return nil
}

func createSuperrewardsIncomeWithTx(tx *sqlx.Tx, income models.Income, transactionID, offerID string) error {
	// insert superrewards income into incomes table
	result, err := addIncome(tx, income)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// update user balance, total_income, referer_total_income
	if err := incrementUserBalance(tx, income.UserID, income.Income, income.RefererIncome); err != nil {
		return err
	}

	// update referer balance
	if _, err := incrementRefererBalance(tx, income.RefererID, income.RefererIncome); err != nil {
		return err
	}

	// insert superrewards offer
	offer := models.SuperrewardsOffer{
		IncomeID:      id,
		UserID:        income.UserID,
		TransactionID: transactionID,
		OfferID:       offerID,
		Amount:        income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `superrewards` (`income_id`, `user_id`, `transaction_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :transaction_id, :offer_id, :amount)", offer)
	return err
}
