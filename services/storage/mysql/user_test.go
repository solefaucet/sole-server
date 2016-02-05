package mysql

import (
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestCreateUser(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When create user", func() {
			err := s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

		Convey("When create user with duplicate email", func() {
			err := s.CreateUser(models.User{Email: "e", BitcoinAddress: ""})

			Convey("Error should be duplicate email", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeDuplicateEmail)
			})
		})

		Convey("When create user with duplicate bitcoin address", func() {
			err := s.CreateUser(models.User{Email: "", BitcoinAddress: "b"})

			Convey("Error should be duplicate bitcoin address", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeDuplicateBitcoinAddress)
			})
		})
	})

	Convey("Given mysql storage with closed connection", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Close()

		Convey("When create user", func() {
			err := s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

			Convey("Error should be unknown", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeUnknown)
			})
		})
	})
}
