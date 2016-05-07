package mysql

import (
	"testing"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRewardRatesByType(t *testing.T) {
	Convey("Given mysql storage with default reward rates", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get reward rates", func() {
			rrs, _ := s.GetRewardRatesByType(models.RewardRateTypeLess)

			Convey("Result set should contains 3 records", func() {
				So(len(rrs), ShouldEqual, 3)
			})
		})
	})

	withClosedConn(t, "When get reward rates", func(s Storage) *errors.Error {
		_, err := s.GetRewardRatesByType(models.RewardRateTypeLess)
		return err
	})
}
