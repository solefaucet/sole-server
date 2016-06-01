package mysql

import (
	"testing"

	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAuthToken(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get auth token", func() {
			_, err := s.GetAuthToken("token")

			Convey("Error should be not found", func() {
				So(err, ShouldEqual, errors.ErrNotFound)
			})
		})
	})

	Convey("Given mysql storage with auth token data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

		Convey("When get auth token", func() {
			authToken, _ := s.GetAuthToken("token")

			Convey("Auth token should be token", func() {
				So(authToken.AuthToken, ShouldEqual, "token")
			})
		})
	})

	withClosedConn(t, "When get auth token", func(s Storage) error {
		_, err := s.GetAuthToken("token")
		return err
	})
}

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
				So(err, ShouldEqual, errors.ErrDuplicatedAuthToken)
			})
		})
	})

	withClosedConn(t, "When create auth token", func(s Storage) error {
		return s.CreateAuthToken(models.AuthToken{AuthToken: "token"})
	})
}

func TestDeleteAuthToken(t *testing.T) {
	Convey("Given mysql storage with auth token data", t, func() {
		s := prepareDatabaseForTesting()
		s.CreateAuthToken(models.AuthToken{AuthToken: "token"})

		Convey("When delete auth token", func() {
			err := s.DeleteAuthToken("token")

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	withClosedConn(t, "When delete auth token", func(s Storage) error {
		return s.DeleteAuthToken("token")
	})
}
