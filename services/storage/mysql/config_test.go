package mysql

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/models"
)

func TestGetLatestConfig(t *testing.T) {
	Convey("Given mysql storage with default config", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get latest config", func() {
			result, err := s.GetLatestConfig()

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Result should be zero value", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					config := actual.(models.Config)
					if config.TotalRewardThreshold == 10 &&
						config.RefererRewardRate == 0.1 {
						return ""
					}
					return fmt.Sprintf("Config %v is not expected", config)
				})
			})
		})
	})

	withClosedConn(t, "When get latest config", func(s Storage) error {
		_, err := s.GetLatestConfig()
		return err
	})
}
