package mysql

import (
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
