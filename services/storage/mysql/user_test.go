package mysql

import (
	"fmt"
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestGetUserByEmail(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get user by email", func() {
			_, err := s.GetUserByEmail("e")

			Convey("Error should be ErrCodeNotFound", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeNotFound)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

		Convey("When get user by email", func() {
			user, _ := s.GetUserByEmail("e")

			Convey("Email should be e, BitcoinAddress should be b", func() {
				So(user, func(actual interface{}, expected ...interface{}) string {
					u := actual.(models.User)
					if u.Email == "e" && u.BitcoinAddress == "b" {
						return ""
					}
					return fmt.Sprintf("User %v is not expected", u)
				})
			})

			Convey("BitcoinAddress should be b", func() {
				So(user.BitcoinAddress, ShouldEqual, "b")
			})
		})
	})

	withClosedConn(t, "When get user by email", func(s Storage) *errors.Error {
		_, err := s.GetUserByEmail("e")
		return err
	})
}
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

	withClosedConn(t, "When create user", func(s Storage) *errors.Error {
		return s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})
	})
}
