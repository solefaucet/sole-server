package v1

import (
	"net/http"
	"testing"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWithdrawalList(t *testing.T) {
	Convey("Given withdrawal list controller", t, func() {
		handler := WithdrawalList(nil, nil)

		Convey("When get withdrawal list with invalid limit", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?limit=3i", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})

		Convey("When get withdrawal list with invalid timestamp", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?since=3i", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})

		Convey("When get withdrawal list without since or until", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})

	Convey("Given withdrawal list controller with errored getWithdrawalSince dependency", t, func() {
		getWithdrawalsSince := mockGetWithdrawalsSince(nil, errors.New(errors.ErrCodeUnknown))
		handler := WithdrawalList(getWithdrawalsSince, nil)

		Convey("When get withdrawal list", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?since=1234567890", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given withdrawal list controller with correct dependencies injected", t, func() {
		getWithdrawalsUntil := mockGetWithdrawalsUntil(nil, nil)
		handler := WithdrawalList(nil, getWithdrawalsUntil)

		Convey("When get withdrawal list", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?until=1234567890", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
