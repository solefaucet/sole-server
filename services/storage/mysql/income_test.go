package mysql

import (
	"fmt"
	"testing"
	"time"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestIncrementUserBalanceByRewardIncome(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		now := time.Now()

		Convey("When increment user balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := incrementUserBalanceByRewardIncome(tx, 0, 0, now)
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})

		Convey("When increment user balance affecting 0 row", func() {
			tx := s.db.MustBegin()
			err := incrementUserBalanceByRewardIncome(tx, 0, 0, now)
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestIncrementRefererBalanceByRewardIncome(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment referer balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			_, err := incrementRefererBalanceByRewardIncome(tx, 0, 0)
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestInsertIntoIncomesTableByRewardIncome(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment referer balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := insertIntoIncomesTableByRewardIncome(tx, models.Income{RefererID: 1})
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestIncrementTotalRewardByRewardIncome(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment total reward with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := incrementTotalRewardByRewardIncome(tx, 10, time.Now())
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestCreateRewardIncome(t *testing.T) {
	Convey("Given mysql storage with two users", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", BitcoinAddress: "b1"})
		s.CreateUser(models.User{Email: "e2", BitcoinAddress: "b2", RefererID: 1})

		Convey("When create reward income", func() {
			err := s.CreateRewardIncome(1, 2, 10, 1, time.Now())

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetIncomes(t *testing.T) {
	withClosedConn(t, "When get incomes", func(s Storage) *errors.Error {
		_, err := s.getIncomes("SELECT * FROM incomes")
		return err
	})
}

func TestGetRewardIncomesSince(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", BitcoinAddress: "b1"})
		rewardedAt := time.Now()
		s.CreateRewardIncome(1, 2, 91, 1, rewardedAt)
		s.CreateRewardIncome(1, 2, 92, 1, rewardedAt)
		s.CreateRewardIncome(1, 2, 93, 1, rewardedAt)

		Convey("When get reward incomes since now", func() {
			result, _ := s.GetRewardIncomesSince(1, time.Now().AddDate(0, 0, -1), 2)

			Convey("Incomes should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					incomes := actual.([]models.Income)
					if len(incomes) == 2 &&
						incomes[0].Income == 91 &&
						incomes[1].Income == 92 {
						return ""
					}
					return fmt.Sprintf("Incomes %v is not expected", incomes)
				})
			})
		})
	})
}

func TestGetRewardIncomesUntil(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", BitcoinAddress: "b1"})
		rewardedAt := time.Now()
		s.CreateRewardIncome(1, 2, 91, 1, rewardedAt)
		s.CreateRewardIncome(1, 2, 92, 1, rewardedAt)
		s.CreateRewardIncome(1, 2, 93, 1, rewardedAt)

		Convey("When get reward incomes until now", func() {
			result, _ := s.GetRewardIncomesUntil(1, time.Now().AddDate(0, 0, 1), 2)

			Convey("Incomes should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					incomes := actual.([]models.Income)
					if len(incomes) == 2 &&
						incomes[0].Income == 93 &&
						incomes[1].Income == 92 {
						return ""
					}
					return fmt.Sprintf("Incomes %v is not expected", incomes)
				})
			})
		})
	})
}
