package mysql

import (
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestCreateAuthToken(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When create auth token", func() {
			err := s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given mysql storage with auth token data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

		Convey("When create auth token with duplicate token", func() {
			err := s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

			Convey("Error should be duplicate auth token", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeDuplicateAuthToken)
			})
		})
	})

	Convey("Given mysql storage with closed connection", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Close()

		Convey("When create auth token", func() {
			err := s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})
		})
	})
}
