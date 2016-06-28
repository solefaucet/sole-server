package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

func TestSuperrewardsCallback(t *testing.T) {
	Convey("Given superrewards callback handler with invalid parameters", t, func() {
		handler := SuperrewardsCallback("", nil, nil, nil, nil, nil)

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})

	Convey("Given superrewards callback handler with invalid signature", t, func() {
		handler := SuperrewardsCallback("secret", nil, nil, nil, nil, nil)
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
		})
	})

	Convey("Given superrewards callback handler with errored getUserByID", t, func() {
		getUserByID := mockGetUserByID(models.User{}, fmt.Errorf(""))
		handler := SuperrewardsCallback("secret", getUserByID, nil, nil, nil, nil)
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496f862b258e8b6b9d3242257"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given superrewards callback handler with non-err-not-found errored getSuperrewardsOfferByID", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getSuperrewardsOfferByID := mockGetSuperrewardsOfferByID(models.SuperrewardsOffer{}, fmt.Errorf(""))
		handler := SuperrewardsCallback("secret", getUserByID, getSuperrewardsOfferByID, nil, nil, nil)
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496f862b258e8b6b9d3242257"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given superrewards callback handler with getSuperrewardsOfferByID returning nil error", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getSuperrewardsOfferByID := mockGetSuperrewardsOfferByID(models.SuperrewardsOffer{}, nil)
		handler := SuperrewardsCallback("secret", getUserByID, getSuperrewardsOfferByID, nil, nil, nil)
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496f862b258e8b6b9d3242257"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})

			Convey("Response body should be 1", func() {
				So(resp.Body.String(), ShouldEqual, "1")
			})
		})
	})

	Convey("Given superrewards callback handler with getSuperrewardsOfferByID returning ErrNotFound and errored createSuperrewardsIncome", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getSuperrewardsOfferByID := mockGetSuperrewardsOfferByID(models.SuperrewardsOffer{}, errors.ErrNotFound)
		getSystemConfig := mockGetSystemConfig(models.Config{})
		createSuperrewardsIncome := mockCreateSuperrewardsIncome(fmt.Errorf(""))
		handler := SuperrewardsCallback("secret", getUserByID, getSuperrewardsOfferByID, getSystemConfig, createSuperrewardsIncome, nil)
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496f862b258e8b6b9d3242257"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given superrewards callback handler with getSuperrewardsOfferByID returning ErrNotFound and correct createSuperrewardsIncome", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		getSuperrewardsOfferByID := mockGetSuperrewardsOfferByID(models.SuperrewardsOffer{}, errors.ErrNotFound)
		getSystemConfig := mockGetSystemConfig(models.Config{})
		createSuperrewardsIncome := mockCreateSuperrewardsIncome(nil)
		handler := SuperrewardsCallback("secret", getUserByID, getSuperrewardsOfferByID, getSystemConfig, createSuperrewardsIncome, func([]byte) {})
		query := "id=id&uid=1&new=13.2&sig=4b2ae6c496f862b258e8b6b9d3242257"

		Convey("When callback", func() {
			route := "/callback"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", route, query), nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
