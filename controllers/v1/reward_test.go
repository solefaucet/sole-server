package v1

import (
	"net/http"
	"testing"
	"time"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestGetReward(t *testing.T) {
	Convey("Given get reward controller with errored getUserByID dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, errors.New(errors.ErrCodeUnknown))
		handler := GetReward(getUserByID, nil, nil, nil, nil, nil)

		Convey("When get reward", func() {
			route := "/incomes/rewards"
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

	Convey("Given get reward controller with not valid last_rewarded", t, func() {
		getUserByID := mockGetUserByID(models.User{RewardedAt: time.Now(), RewardInterval: 5}, nil)
		handler := GetReward(getUserByID, nil, nil, nil, nil, nil)

		Convey("When get reward", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
		})
	})

	Convey("Given get reward controller with errored createRewardIncome dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getLatestTotalReward := mockGetLatestTotalReward(models.TotalReward{CreatedAt: time.Now().UTC(), Total: 11})
		getSystemConfig := mockGetSystemConfig(models.Config{TotalRewardThreshold: 10, RefererRewardRate: 0.1})
		getRewardRatesByType := mockGetRewardRatesByType([]models.RewardRate{
			{Weight: 1, Min: 1, Max: 10},
			{Weight: 2, Min: 11, Max: 20},
			{Weight: 3, Min: 21, Max: 30},
		})
		getBitcoinPrice := mockGetBitcoinPrice(40000000)
		createRewardIncome := mockCreateRewardIncome(errors.New(errors.ErrCodeUnknown))
		handler := GetReward(getUserByID, getLatestTotalReward, getSystemConfig, getRewardRatesByType, getBitcoinPrice, createRewardIncome)

		Convey("When get reward", func() {
			route := "/incomes/rewards"
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

	Convey("Given get reward controller with everything correctly configured", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getLatestTotalReward := mockGetLatestTotalReward(models.TotalReward{CreatedAt: time.Now().UTC(), Total: 11})
		getSystemConfig := mockGetSystemConfig(models.Config{TotalRewardThreshold: 10, RefererRewardRate: 0.1})
		getRewardRatesByType := mockGetRewardRatesByType([]models.RewardRate{
			{Weight: 1, Min: 1, Max: 10},
			{Weight: 2, Min: 11, Max: 20},
			{Weight: 3, Min: 21, Max: 30},
		})
		getBitcoinPrice := mockGetBitcoinPrice(40000000)
		createRewardIncome := mockCreateRewardIncome(nil)
		handler := GetReward(getUserByID, getLatestTotalReward, getSystemConfig, getRewardRatesByType, getBitcoinPrice, createRewardIncome)

		Convey("When get reward", func() {
			route := "/incomes/rewards"
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

func TestRewardList(t *testing.T) {
	Convey("Given reward list controller", t, func() {
		handler := RewardList(nil, nil)

		Convey("When get reward list with invalid limit", func() {
			route := "/incomes/rewards"
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

		Convey("When get reward list with invalid timestamp", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?timestamp=3i", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})

		Convey("When get reward list with invalid type", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?timestamp=3&type=wow", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})

	Convey("Given reward list controller with errored getRewardIncomesSince dependency", t, func() {
		getRewardIncomesSince := mockGetRewardIncomesSince(nil, errors.New(errors.ErrCodeUnknown))
		handler := RewardList(getRewardIncomesSince, nil)

		Convey("When get reward list", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?timestamp=3&type=since", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given reward list controller with correct dependencies injected", t, func() {
		getRewardIncomesUntil := mockGetRewardIncomesUntil(nil, nil)
		handler := RewardList(nil, getRewardIncomesUntil)

		Convey("When get reward list", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?timestamp=3&type=until", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
