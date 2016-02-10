package v1

import (
	"net/http"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestRequestVerifyEmail(t *testing.T) {
	Convey("Given request verify email controller with errored getUserByID dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, errors.New(errors.ErrCodeUnknown))
		handler := RequestVerifyEmail(getUserByID, nil, nil)

		Convey("When request verify email", func() {
			route := "/sessions"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given request verify email controller with errored upsertSession dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		upsertSession := mockUpsertSession(errors.New(errors.ErrCodeUnknown))
		handler := RequestVerifyEmail(getUserByID, upsertSession, nil)

		Convey("When request verify email", func() {
			route := "/sessions"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given request verify email controller with errored sendEmail dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		upsertSession := mockUpsertSession(nil)
		sendEmail := mockSendEmail(errors.New(errors.ErrCodeUnknown))
		handler := RequestVerifyEmail(getUserByID, upsertSession, sendEmail)

		Convey("When request verify email", func() {
			route := "/sessions"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given request verify email controller with correct dependencies injected", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		upsertSession := mockUpsertSession(nil)
		sendEmail := mockSendEmail(nil)
		handler := RequestVerifyEmail(getUserByID, upsertSession, sendEmail)

		Convey("When request verify email", func() {
			route := "/sessions"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
