package v1

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/freeusd/solebtc/models"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetReward(t *testing.T) {
	Convey("Given get reward controller with errored getUserByID dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, fmt.Errorf(""))
		handler := GetReward(getUserByID, nil, nil, nil, nil, nil, nil)

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
		handler := GetReward(getUserByID, nil, nil, nil, nil, nil, nil)

		Convey("When get reward", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 429", func() {
				So(resp.Code, ShouldEqual, statusCodeTooManyRequests)
			})
		})
	})

	Convey("Given get reward controller with errored createRewardIncome dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getLatestTotalReward := mockGetLatestTotalReward(models.TotalReward{CreatedAt: time.Now(), Total: 11})
		getSystemConfig := mockGetSystemConfig(models.Config{TotalRewardThreshold: 10, RefererRewardRate: 10})
		getRewardRatesByType := mockGetRewardRatesByType([]models.RewardRate{
			{Weight: 1, Min: 1, Max: 10},
			{Weight: 2, Min: 11, Max: 20},
			{Weight: 3, Min: 21, Max: 30},
		})
		createRewardIncome := mockCreateRewardIncome(fmt.Errorf(""))
		handler := GetReward(getUserByID, getLatestTotalReward, getSystemConfig, getRewardRatesByType, createRewardIncome, nil, nil)

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
		getSystemConfig := mockGetSystemConfig(models.Config{TotalRewardThreshold: 10, RefererRewardRate: 10})
		getRewardRatesByType := mockGetRewardRatesByType([]models.RewardRate{
			{Weight: 1, Min: 1, Max: 10},
			{Weight: 2, Min: 11, Max: 20},
			{Weight: 3, Min: 21, Max: 30},
		})
		createRewardIncome := mockCreateRewardIncome(nil)
		insertIncome := mockInsertIncome()
		broadcast := mockBroadcast()
		handler := GetReward(getUserByID,
			getLatestTotalReward,
			getSystemConfig,
			getRewardRatesByType,
			createRewardIncome,
			insertIncome,
			broadcast,
		)

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
	Convey("Given reward list handler", t, func() {
		handler := RewardList(mockGetRewardIncomes(nil, fmt.Errorf("")), nil)

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

		Convey("When get reward list with errored handler", func() {
			route := "/incomes/rewards"
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

	Convey("Given reward list handler", t, func() {
		handler := RewardList(mockGetRewardIncomes(nil, nil), func(int64) (int64, error) { return 0, nil })

		Convey("When get reward list", func() {
			route := "/incomes/rewards"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
