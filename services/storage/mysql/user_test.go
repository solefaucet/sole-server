package mysql

import (
	"fmt"
	"testing"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUserByID(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get user by id", func() {
			_, err := s.GetUserByID(1)

			Convey("Error should be ErrNotFound", func() {
				So(err, ShouldEqual, errors.ErrNotFound)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", Address: "b"})

		Convey("When get user by id", func() {
			user, _ := s.GetUserByID(1)

			Convey("ID should be 1, Email should be e, Address should be b", func() {
				So(user, func(actual interface{}, expected ...interface{}) string {
					u := actual.(models.User)
					if u.ID == 1 &&
						u.Email == "e" &&
						u.Address == "b" {
						return ""
					}
					return fmt.Sprintf("User %v is not expected", u)
				})
			})
		})
	})

	withClosedConn(t, "When get user by id", func(s Storage) error {
		_, err := s.GetUserByID(1)
		return err
	})
}

func TestGetUserByEmail(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get user by email", func() {
			_, err := s.GetUserByEmail("e")

			Convey("Error should be ErrNotFound", func() {
				So(err, ShouldEqual, errors.ErrNotFound)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", Address: "b"})

		Convey("When get user by email", func() {
			user, _ := s.GetUserByEmail("e")

			Convey("Email should be e, Address should be b", func() {
				So(user, func(actual interface{}, expected ...interface{}) string {
					u := actual.(models.User)
					if u.Email == "e" && u.Address == "b" {
						return ""
					}
					return fmt.Sprintf("User %v is not expected", u)
				})
			})
		})
	})

	withClosedConn(t, "When get user by email", func(s Storage) error {
		_, err := s.GetUserByEmail("e")
		return err
	})
}

func TestCreateUser(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When create user", func() {
			err := s.CreateUser(models.User{Email: "e", Address: "b"})

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", Address: "b"})

		Convey("When create user with duplicate email", func() {
			err := s.CreateUser(models.User{Email: "e", Address: ""})

			Convey("Error should be duplicate email", func() {
				So(err, ShouldEqual, errors.ErrDuplicatedEmail)
			})
		})

		Convey("When create user with duplicate address", func() {
			err := s.CreateUser(models.User{Email: "", Address: "b"})

			Convey("Error should be duplicate address", func() {
				So(err, ShouldEqual, errors.ErrDuplicatedAddress)
			})
		})
	})

	withClosedConn(t, "When create user", func(s Storage) error {
		return s.CreateUser(models.User{Email: "e", Address: "b"})
	})
}

func TestUpdateUserStatus(t *testing.T) {
	Convey("Given mysql storage with user data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e", Address: "b"})

		Convey("When update user's status", func() {
			err := s.UpdateUserStatus(1, models.UserStatusVerified)
			user, _ := s.GetUserByID(1)

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("New user status should be verified", func() {
				So(user.Status, ShouldEqual, models.UserStatusVerified)
			})
		})
	})

	withClosedConn(t, "When update user", func(s Storage) error {
		return s.UpdateUserStatus(0, "")
	})
}

func TestGetReferees(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateUser(models.User{Email: "e1", Address: "b1"})
		s.CreateUser(models.User{RefererID: 1, Email: "e2", Address: "b2"})
		s.CreateUser(models.User{RefererID: 1, Email: "e3", Address: "b3"})

		Convey("When get referees until 1", func() {
			result, _ := s.GetReferees(1, 1, 1)

			Convey("Users should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					users := actual.([]models.User)
					if len(users) == 1 &&
						users[0].Email == "e2" {
						return ""
					}
					return fmt.Sprintf("Users %v is not expected", result)
				})
			})
		})
	})
}

func TestGetWithdrawableUser(t *testing.T) {
	Convey("Given mysql storage", t, func() {
		s := prepareDatabaseForTesting()
		s.db.Exec("INSERT INTO users(email, address, status, balance, min_withdrawal_amount) VALUES('e1', 'b1', 'verified', 10, 5)")
		s.db.Exec("INSERT INTO users(email, address, status, balance, min_withdrawal_amount) VALUES('e2', 'b2', 'verified', 5, 10)")

		Convey("When get withdrawable users", func() {
			result, _ := s.GetWithdrawableUsers()

			Convey("Users should equal", func() {
				So(result, func(actual interface{}, expected ...interface{}) string {
					users := actual.([]models.User)
					if len(users) == 1 && users[0].Email == "e1" {
						return ""
					}
					return fmt.Sprintf("Users %v is not expected", result)
				})
			})
		})
	})
}
