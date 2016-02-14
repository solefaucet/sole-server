package mysql

import (
	"fmt"
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestGetUserByID(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get user by id", func() {
			_, err := s.GetUserByID(1)

			Convey("Error should be ErrCodeNotFound", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeNotFound)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

		Convey("When get user by id", func() {
			user, _ := s.GetUserByID(1)

			Convey("ID should be 1, Email should be e, BitcoinAddress should be b", func() {
				So(user, func(actual interface{}, expected ...interface{}) string {
					u := actual.(models.User)
					if u.ID == 1 &&
						u.Email == "e" &&
						u.BitcoinAddress == "b" {
						return ""
					}
					return fmt.Sprintf("User %v is not expected", u)
				})
			})
		})
	})

	withClosedConn(t, "When get user by id", func(s Storage) *errors.Error {
		_, err := s.GetUserByID(1)
		return err
	})
}

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

func TestUpdateUser(t *testing.T) {
	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", BitcoinAddress: "b"})

		Convey("When update user", func() {
			err := s.UpdateUser(models.User{ID: 1, Status: models.UserStatusVerified})
			user, _ := s.GetUserByID(1)

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("New user status should be verified", func() {
				So(user.Status, ShouldEqual, models.UserStatusVerified)
			})
		})
	})

	withClosedConn(t, "When update user", func(s Storage) *errors.Error {
		return s.UpdateUser(models.User{})
	})
}

func TestGetRefereesSince(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", BitcoinAddress: "b1"})
		s.CreateUser(models.User{RefererID: 1, Email: "e2", BitcoinAddress: "b2"})
		s.CreateUser(models.User{RefererID: 1, Email: "e3", BitcoinAddress: "b3"})

		Convey("When get referees since 1", func() {
			result, _ := s.GetRefereesSince(1, 1, 2)

			Convey("Users should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					users := actual.([]models.User)
					if len(users) == 2 &&
						users[0].Email == "e2" &&
						users[1].Email == "e3" {
						return ""
					}
					return fmt.Sprintf("Users %v is not expected", result)
				})
			})
		})
	})
}

func TestGetRefereesUntil(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", BitcoinAddress: "b1"})
		s.CreateUser(models.User{RefererID: 1, Email: "e2", BitcoinAddress: "b2"})
		s.CreateUser(models.User{RefererID: 1, Email: "e3", BitcoinAddress: "b3"})

		Convey("When get referees until 1", func() {
			result, _ := s.GetRefereesUntil(1, 10, 1)

			Convey("Users should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					users := actual.([]models.User)
					if len(users) == 1 &&
						users[0].Email == "e3" {
						return ""
					}
					return fmt.Sprintf("Users %v is not expected", result)
				})
			})
		})
	})
}
