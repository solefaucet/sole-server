package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/solefaucet/sole-server/models"
)

// ChargebackIncome set income status to chargeback
func (s Storage) ChargebackIncome(incomeID int64) error {
	sql := "UPDATE `incomes` SET `status` = ? WHERE `id` = ?"
	args := []interface{}{models.IncomeStatusChargeback, incomeID}
	_, err := s.db.Exec(sql, args...)
	return err
}

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
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `incomes` WHERE `user_id` = ? AND `type` = ?", userID, models.IncomeTypeReward).Scan(&count)
	return count, err
}

// GetOfferwallIncomes get user's offerwall incomes
func (s Storage) GetOfferwallIncomes(userID int64, limit, offset int64) ([]models.Income, error) {
	rawSQL := "SELECT * FROM incomes WHERE `user_id` = ? AND `type` != ? ORDER BY `id` DESC LIMIT ? OFFSET ?"
	args := []interface{}{userID, models.IncomeTypeReward, limit, offset}
	incomes := []models.Income{}
	err := s.selects(&incomes, rawSQL, args...)
	return incomes, err
}

// GetNumberOfOfferwallIncomes gets number of user's offerwall incomes
func (s Storage) GetNumberOfOfferwallIncomes(userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `incomes` WHERE `user_id` = ? AND `type` != ?", userID, models.IncomeTypeReward).Scan(&count)
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

	_, rowAffected, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}
	if rowAffected == 1 {
		totalReward += income.RefererIncome
	}

	// update user rewarded_at
	if _, err := tx.Exec("UPDATE users SET `rewarded_at` = ? WHERE `id` = ?", now, income.UserID); err != nil {
		return err
	}

	// update total reward
	if err := incrementTotalReward(tx, totalReward, now); err != nil {
		return err
	}

	return nil
}

// GetNumberOfOfferwowEvents gets number of offerwow event
func (s Storage) GetNumberOfOfferwowEvents(eventID string) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `offerwow` WHERE `event_id` = ?", eventID).Scan(&count)
	return count, err
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
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
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

// GetNumberOfSuperrewardsOffers gets number of superrewards offers
func (s Storage) GetNumberOfSuperrewardsOffers(transactionID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `superrewards` WHERE `transaction_id` = ? AND `user_id` = ?", transactionID, userID).Scan(&count)
	return count, err
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
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
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

// GetNumberOfKiwiwallOffers gets number of kiwiwall offers
func (s Storage) GetNumberOfKiwiwallOffers(transactionID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `kiwiwall` WHERE `transaction_id` = ? AND `user_id` = ?", transactionID, userID).Scan(&count)
	return count, err
}

// CreateKiwiwallIncome creates a new kiwiwall type income
func (s Storage) CreateKiwiwallIncome(income models.Income, transactionID, offerID string) error {
	tx := s.db.MustBegin()

	if err := createKiwiwallIncomeWithTx(tx, income, transactionID, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create kiwiwall income commit transaction error: %v", err)
	}

	return nil
}

func createKiwiwallIncomeWithTx(tx *sqlx.Tx, income models.Income, transactionID, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert kiwiwall offer
	offer := models.KiwiwallOffer{
		IncomeID:      id,
		UserID:        income.UserID,
		TransactionID: transactionID,
		OfferID:       offerID,
		Amount:        income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `kiwiwall` (`income_id`, `user_id`, `transaction_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :transaction_id, :offer_id, :amount)", offer)
	return err
}

// GetAdscendMediaOffer returns AdscendMediaOffer
func (s Storage) GetAdscendMediaOffer(transactionID string, userID int64) (*models.AdscendMedia, error) {
	dest := &models.AdscendMedia{}
	query := "SELECT * FROM `adscend_media` WHERE `transaction_id` = ? AND `user_id` = ?"
	args := []interface{}{transactionID, userID}
	err := s.db.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return dest, err
}

// CreateAdscendMediaIncome creates a new adscend media type income
func (s Storage) CreateAdscendMediaIncome(income models.Income, transactionID, offerID string) error {
	tx := s.db.MustBegin()

	if err := createAdscendMediaIncomeWithTx(tx, income, transactionID, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create adscend media income commit transaction error: %v", err)
	}

	return nil
}

func createAdscendMediaIncomeWithTx(tx *sqlx.Tx, income models.Income, transactionID, offerID string) error {
	// insert income into incomes table
	result, err := addIncome(tx, income)
	if err != nil {
		return err
	}
	incomeID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// insert adscend media offer
	offer := models.AdscendMedia{
		IncomeID:      incomeID,
		UserID:        income.UserID,
		TransactionID: transactionID,
		OfferID:       offerID,
		Amount:        income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `adscend_media` (`income_id`, `user_id`, `transaction_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :transaction_id, :offer_id, :amount)", offer)
	return err
}

// GetNumberOfAdgateMediaOffers gets number of adgate media offers
func (s Storage) GetNumberOfAdgateMediaOffers(transactionID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `adgate_media` WHERE `transaction_id` = ? AND `user_id` = ?", transactionID, userID).Scan(&count)
	return count, err
}

// CreateAdgateMediaIncome creates a new adgate media type income
func (s Storage) CreateAdgateMediaIncome(income models.Income, transactionID, offerID string) error {
	tx := s.db.MustBegin()

	if err := createAdgateMediaIncomeWithTx(tx, income, transactionID, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create adgate media income commit transaction error: %v", err)
	}

	return nil
}

func createAdgateMediaIncomeWithTx(tx *sqlx.Tx, income models.Income, transactionID, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert adgate media offer
	offer := models.AdgateMedia{
		IncomeID:      id,
		UserID:        income.UserID,
		TransactionID: transactionID,
		OfferID:       offerID,
		Amount:        income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `adgate_media` (`income_id`, `user_id`, `transaction_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :transaction_id, :offer_id, :amount)", offer)
	return err
}

// GetNumberOfOffertoroOffers gets number of offertoro offers
func (s Storage) GetNumberOfOffertoroOffers(transactionID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `offertoro` WHERE `transaction_id` = ? AND `user_id` = ?", transactionID, userID).Scan(&count)
	return count, err
}

// CreateOffertoroIncome creates a new offertoro type income
func (s Storage) CreateOffertoroIncome(income models.Income, transactionID, offerID string) error {
	tx := s.db.MustBegin()

	if err := createOffertoroIncomeWithTx(tx, income, transactionID, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create offertoro income commit transaction error: %v", err)
	}

	return nil
}

func createOffertoroIncomeWithTx(tx *sqlx.Tx, income models.Income, transactionID, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert offertoro offer
	offer := models.Offertoro{
		IncomeID:      id,
		UserID:        income.UserID,
		TransactionID: transactionID,
		OfferID:       offerID,
		Amount:        income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `offertoro` (`income_id`, `user_id`, `transaction_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :transaction_id, :offer_id, :amount)", offer)
	return err
}

// GetNumberOfPersonalyOffers gets number of personaly offers
func (s Storage) GetNumberOfPersonalyOffers(offerID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `personaly` WHERE `offer_id` = ? AND `user_id` = ?", offerID, userID).Scan(&count)
	return count, err
}

// CreatePersonalyIncome creates a new personaly type income
func (s Storage) CreatePersonalyIncome(income models.Income, offerID string) error {
	tx := s.db.MustBegin()

	if err := createPersonalyIncomeWithTx(tx, income, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create personaly income commit transaction error: %v", err)
	}

	return nil
}

func createPersonalyIncomeWithTx(tx *sqlx.Tx, income models.Income, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert personaly offer
	offer := models.PersonalyOffer{
		IncomeID: id,
		UserID:   income.UserID,
		OfferID:  offerID,
		Amount:   income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `personaly` (`income_id`, `user_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :offer_id, :amount)", offer)
	return err
}

// GetNumberOfTrialpayOffers gets number of trialpay offers
func (s Storage) GetNumberOfTrialpayOffers(offerID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `trialpay` WHERE `offer_id` = ? AND `user_id` = ?", offerID, userID).Scan(&count)
	return count, err
}

// CreateTrialpayIncome creates a new trialpay type income
func (s Storage) CreateTrialpayIncome(income models.Income, offerID string) error {
	tx := s.db.MustBegin()

	if err := createTrialpayIncomeWithTx(tx, income, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create trialpay income commit transaction error: %v", err)
	}

	return nil
}

func createTrialpayIncomeWithTx(tx *sqlx.Tx, income models.Income, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert trialpay offer
	offer := models.TrialpayOffer{
		IncomeID: id,
		UserID:   income.UserID,
		OfferID:  offerID,
		Amount:   income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `trialpay` (`income_id`, `user_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :offer_id, :amount)", offer)
	return err
}

// GetNumberOfClixwallOffers gets number of clixwall offers
func (s Storage) GetNumberOfClixwallOffers(offerID string, userID int64) (int64, error) {
	var count int64
	err := s.db.QueryRowx("SELECT COUNT(*) FROM `clixwalls` WHERE `offer_id` = ? AND `user_id` = ?", offerID, userID).Scan(&count)
	return count, err
}

// CreateClixwallIncome creates a new clixwall type income
func (s Storage) CreateClixwallIncome(income models.Income, offerID string) error {
	tx := s.db.MustBegin()

	if err := createClixwallIncomeWithTx(tx, income, offerID); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create clixwall income commit transaction error: %v", err)
	}

	return nil
}

func createClixwallIncomeWithTx(tx *sqlx.Tx, income models.Income, offerID string) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert clixwall offer
	offer := models.ClixwallOffer{
		IncomeID: id,
		UserID:   income.UserID,
		OfferID:  offerID,
		Amount:   income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `clixwalls` (`income_id`, `user_id`, `offer_id`, `amount`) VALUE (:income_id, :user_id, :offer_id, :amount)", offer)
	return err
}

// CreatePtcwallIncome creates a new ptcwall type income
func (s Storage) CreatePtcwallIncome(income models.Income) error {
	tx := s.db.MustBegin()

	if err := createPtcwallIncomeWithTx(tx, income); err != nil {
		tx.Rollback()
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("create clixwall income commit transaction error: %v", err)
	}

	return nil
}

func createPtcwallIncomeWithTx(tx *sqlx.Tx, income models.Income) error {
	id, _, err := commonBatchOperation(tx, income)
	if err != nil {
		return err
	}

	// insert ptcwall offer
	offer := models.PtcwallOffer{
		IncomeID: id,
		UserID:   income.UserID,
		Amount:   income.Income,
	}
	_, err = tx.NamedExec("INSERT INTO `ptcwalls` (`income_id`, `user_id`, `amount`) VALUE (:income_id, :user_id, :amount)", offer)
	return err
}

// add income, update user, update referer
func commonBatchOperation(tx *sqlx.Tx, income models.Income) (incomeID, updateRefererBalanceRowsAffected int64, err error) {
	// insert income into incomes table
	result, err := addIncome(tx, income)
	if err != nil {
		return
	}
	incomeID, err = result.LastInsertId()
	if err != nil {
		return
	}

	// update user balance, total_income, referer_total_income
	if err = incrementUserBalance(tx, income.UserID, income.Income, income.RefererIncome); err != nil {
		return
	}

	// update referer balance
	updateRefererBalanceRowsAffected, err = incrementRefererBalance(tx, income.RefererID, income.RefererIncome)
	if _, err = incrementRefererBalance(tx, income.RefererID, income.RefererIncome); err != nil {
		return
	}

	return
}

// insert reward income into incomes table
func addIncome(tx *sqlx.Tx, income models.Income) (sql.Result, error) {
	sql := "INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`, `status`) VALUES (:user_id, :referer_id, :type, :income, :referer_income, :status)"

	if income.Status == "" {
		sql = "INSERT INTO incomes (`user_id`, `referer_id`, `type`, `income`, `referer_income`) VALUES (:user_id, :referer_id, :type, :income, :referer_income)"
	}

	result, err := tx.NamedExec(sql, income)
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
