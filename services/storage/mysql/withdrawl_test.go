package mysql

import (
	"testing"
	"time"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestDeductUserBalanceBy(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When deduct user balance with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := deductUserBalanceBy(tx, 0, 0)
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})

		Convey("When deduct user balance affecting 0 row", func() {
			tx := s.db.MustBegin()
			err := deductUserBalanceBy(tx, 0, 0)
			Convey("Error should be insufficient balance", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeInsufficientBalance)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestInsertWithdrawl(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When insert withdrawl with commited transaction", func() {
			tx := s.db.MustBegin()
			tx.Commit()
			err := insertWithdrawl(tx, 0, "", 0)
			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})

			Reset(func() { tx.Rollback() })
		})
	})
}

func TestCreateWithdrawl(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})
		s.CreateRewardIncome(models.Income{UserID: 1, Income: 10}, time.Now())

		Convey("When create withdrawl", func() {
			err := s.CreateWithdrawl(models.Withdrawl{
				UserID:         1,
				BitcoinAddress: "b",
				Amount:         5,
			})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
