package mysql

import (
	"fmt"
	"testing"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetSessionByToken(t *testing.T) {
	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When get session", func() {
			_, err := s.GetSessionByToken("token")

			Convey("Error should be not found", func() {
				So(err, ShouldEqual, errors.ErrNotFound)
			})
		})
	})

	Convey("Given mysql storage with session data", t, func() {
		s := prepareDatabaseForTesting()
		s.UpsertSession(models.Session{
			UserID: 1,
			Token:  "token",
			Type:   "verify-email",
		})

		Convey("When get session", func() {
			session, _ := s.GetSessionByToken("token")

			Convey("Session should equal", func() {
				So(session, func(actual interface{}, expected ...interface{}) string {
					s := actual.(models.Session)
					if s.UserID == 1 &&
						s.Token == "token" &&
						s.Type == "verify-email" {
						return ""
					}
					return fmt.Sprintf("Session %v is not expected", s)
				})
			})
		})
	})

	withClosedConn(t, "When get session", func(s Storage) error {
		_, err := s.GetSessionByToken("token")
		return err
	})
}

func TestUpsertSession(t *testing.T) {
	sess := models.Session{
		UserID: 1,
		Token:  "token",
		Type:   "verify-email",
	}

	Convey("Given empty mysql storage", t, func() {
		s := prepareDatabaseForTesting()

		Convey("When upsert session", func() {
			err := s.UpsertSession(sess)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given mysql storage with session data", t, func() {
		s := prepareDatabaseForTesting()
		s.UpsertSession(sess)

		Convey("When upsert session with duplicate token", func() {
			err := s.UpsertSession(sess)

			Convey("Error should also be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	withClosedConn(t, "When upsert session", func(s Storage) error {
		return s.UpsertSession(sess)
	})
}
