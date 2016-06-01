package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/solefaucet/solebtc/models"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWithdrawalList(t *testing.T) {
	Convey("Given withdrawal list controller", t, func() {
		handler := WithdrawalList(mockGetWithdrawals(nil, fmt.Errorf("")), nil, nil)

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

		Convey("When get withdrawal list", func() {
			route := "/withdrawals"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given withdrawal list controller with correct dependencies injected", t, func() {
		getWithdrawals := mockGetWithdrawals([]models.Withdrawal{{}}, nil)
		handler := WithdrawalList(getWithdrawals, func(int64) (int64, error) { return 0, nil }, func(tx string) string { return tx })

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
