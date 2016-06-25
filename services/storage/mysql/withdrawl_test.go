package mysql

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

func TestDeductUserBalanceBy(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When deduct user balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := deductUserBalanceBy(tx, 0, 0)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})

		Convey("When deduct user balance affecting 0 row", func() {
			tx := s.db.MustBegin()
			err := deductUserBalanceBy(tx, 0, 0)
			Convey("Error should be insufficient balance", func() {
				So(err, ShouldEqual, errors.ErrInsufficientBalance)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestInsertWithdrawal(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When insert withdrawal with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := insertWithdrawal(tx, 0, "", 0)
			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestCreateWithdrawal(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", Address: "b"})
		s.CreateRewardIncome(models.Income{UserID: 1, Income: 10}, time.Now())

		Convey("When create withdrawal", func() {
			err := s.CreateWithdrawal(models.Withdrawal{
				UserID:  1,
				Address: "b",
				Amount:  5,
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetWithdrawals(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.db.MustExec("INSERT INTO `users` (email, address, balance) VALUES(?, ?, ?);", "e", "b", 8388607)
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 1})
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 2})
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 3})

		Convey("When get withdrawals until now", func() {
			result, _ := s.GetWithdrawals(1, 2, 1)

			Convey("Withdrawals should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					withdrawals := actual.([]models.Withdrawal)
					if len(withdrawals) == 2 &&
						withdrawals[0].Amount == 2 &&
						withdrawals[1].Amount == 1 {
						return ""
					}
					return fmt.Sprintf("Withdrawals %v is not expected", withdrawals)
				})
			})
		})
	})
}

func TestGetUnprocessedWithdrawals(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.db.MustExec("INSERT INTO `users` (email, address, balance) VALUES(?, ?, ?);", "e", "b", 8388607)
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 1})
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 2})
		s.CreateWithdrawal(models.Withdrawal{UserID: 1, Address: "b", Amount: 3})

		Convey("When get withdrawals until now", func() {
			result, _ := s.GetPendingWithdrawals()

			Convey("Withdrawals should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					withdrawals := actual.([]models.Withdrawal)
					if withdrawals[0].Amount == 1 &&
						withdrawals[1].Amount == 2 &&
						withdrawals[2].Amount == 3 {
						return ""
					}
					return fmt.Sprintf("Withdrawals %v is not expected", withdrawals)
				})
			})
		})
	})
}
