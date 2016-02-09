package mysql

import (
	"fmt"
	"testing"
	"time"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestGetSortedTotalRewards(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get total rewards", func() {
			trs, err := s.GetSortedTotalRewards()

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result set should be empty", func() {
				So(trs, ShouldBeEmpty)
			})
		})
	})

	Convey("Given mysql storage with one total reward data", t, func() {
		s := prepareDatabaseForTesting()
		now := time.Now()
		s.IncrementTotalReward(now, 1)
		s.IncrementTotalReward(now, 1)

		Convey("When get total rewards", func() {
			trs, _ := s.GetSortedTotalRewards()

			Convey("Result set should be equal", func() {
				So(trs, func(actual interface{}, expected ...interface{}) string {
					result := actual.([]models.TotalReward)
					if len(result) == 1 &&
						result[0].CreatedAt.YearDay() == now.YearDay() &&
						result[0].Total == 2 {
						return ""
					}
					return fmt.Sprintf("Result set %v is not expected", result)
				})
			})
		})
	})

	Convey("Given mysql storage with two total reward data", t, func() {
		s := prepareDatabaseForTesting()
		now := time.Now()
		tmr := now.AddDate(0, 0, 1)
		s.IncrementTotalReward(now, 10)
		s.IncrementTotalReward(tmr, 1)

		Convey("When get total rewards", func() {
			trs, _ := s.GetSortedTotalRewards()

			Convey("Result set should be equal", func() {
				So(trs, func(actual interface{}, expected ...interface{}) string {
					result := actual.([]models.TotalReward)
					if len(result) == 2 &&
						result[0].Total == 1 &&
						result[0].CreatedAt.YearDay() == tmr.YearDay() &&
						result[1].Total == 10 &&
						result[1].CreatedAt.YearDay() == now.YearDay() {
						return ""
					}
					return fmt.Sprintf("Result set %v is not expected", result)
				})
			})
		})
	})

	withClosedConn(t, "When total rewards", func(s Storage) *errors.Error {
		_, err := s.GetSortedTotalRewards()
		return err
	})
}

func TestIncrementTotalReward(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When increment total reward", func() {
			err := s.IncrementTotalReward(time.Now(), 1)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	withClosedConn(t, "When increment total reward", func(s Storage) *errors.Error {
		return s.IncrementTotalReward(time.Now(), 1)
	})
}
