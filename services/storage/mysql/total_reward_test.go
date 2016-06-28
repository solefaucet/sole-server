package mysql

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/models"
)

func TestGetLatestTotalReward(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get latest total reward", func() {
			result, err := s.GetLatestTotalReward()

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result should be zero value", func() {
				So(result, ShouldResemble, models.TotalReward{})
			})
		})
	})

	Convey("Given mysql storage with one total reward data", t, func() {
		s := prepareDatabaseForTesting()
		now := time.Now()
		s.IncrementTotalReward(now, 1)
		s.IncrementTotalReward(now, 1)

		Convey("When get latest total reward", func() {
			r, _ := s.GetLatestTotalReward()

			Convey("Result should be equal", func() {
				So(r, func(actual interface{}, expected ...interface{}) string {
					result := actual.(models.TotalReward)
					if result.IsSameDay(now) && result.Total == 2 {
						return ""
					}
					return fmt.Sprintf("Result %v is not expected", result)
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

		Convey("When get latest total reward", func() {
			r, _ := s.GetLatestTotalReward()

			Convey("Result should be equal", func() {
				So(r, func(actual interface{}, expected ...interface{}) string {
					result := actual.(models.TotalReward)
					if result.IsSameDay(tmr) && result.Total == 1 {
						return ""
					}
					return fmt.Sprintf("Result %v is not expected", result)
				})
			})
		})
	})

	withClosedConn(t, "When get latest total rewards", func(s Storage) error {
		_, err := s.GetLatestTotalReward()
		return err
	})
}
