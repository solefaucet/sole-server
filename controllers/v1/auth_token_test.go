package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func TestLogin(t *testing.T) {
	requestDataJSON := func(email string) []byte {
		raw, _ := json.Marshal(map[string]interface{}{
			"email": email,
		})
		return raw
	}

	testdata := []struct {
		when            string
		requestData     []byte
		code            int
		getUserByEmail  dependencyGetUserByEmail
		createAuthToken dependencyCreateAuthToken
	}{
		{
			"invalid json data",
			[]byte("huhu"),
			400,
			nil,
			nil,
		},
		{
			"invalid email",
			requestDataJSON(invalidEmail),
			400,
			nil,
			nil,
		},
		{
			"banned user",
			requestDataJSON(validEmail),
			403,
			mockGetUserByEmail(models.User{Status: models.UserStatusBanned}, nil),
			nil,
		},
		{
			"non existing email",
			requestDataJSON(validEmail),
			404,
			mockGetUserByEmail(models.User{}, errors.New(errors.ErrCodeNotFound)),
			nil,
		},
		{
			"valid email, unknown error",
			requestDataJSON(validEmail),
			500,
			mockGetUserByEmail(models.User{}, errors.New(errors.ErrCodeUnknown)),
			nil,
		},
		{
			"valid existing email, but unknown error",
			requestDataJSON(validEmail),
			500,
			mockGetUserByEmail(models.User{}, nil),
			mockCreateAuthToken(errors.New(errors.ErrCodeUnknown)),
		},
		{
			"valid existing email",
			requestDataJSON(validEmail),
			201,
			mockGetUserByEmail(models.User{}, nil),
			mockCreateAuthToken(nil),
		},
	}

	for _, v := range testdata {
		Convey("Given Login controller", t, func() {
			handler := Login(v.getUserByEmail, v.createAuthToken)

			Convey(fmt.Sprintf("When request with %s", v.when), func() {
				route := "/auth_tokens"
				_, resp, r := gin.CreateTestContext()
				r.POST(route, handler)
				req, _ := http.NewRequest("POST", route, bytes.NewBuffer(v.requestData))
				r.ServeHTTP(resp, req)

				Convey(fmt.Sprintf("Response code should be equal to %d", v.code), func() {
					So(resp.Code, ShouldEqual, v.code)
				})
			})
		})
	}
}

func TestLogout(t *testing.T) {
	Convey("Given Logout controller with errored logout dependency", t, func() {
		handler := Logout(mockDeleteAuthToken(errors.New(errors.ErrCodeUnknown)))

		Convey("When logout", func() {
			route := "/auth_tokens"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.DELETE(route, handler)
			req, _ := http.NewRequest("DELETE", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given Logout controller", t, func() {
		handler := Logout(mockDeleteAuthToken(nil))

		Convey("When logout", func() {
			route := "/auth_tokens"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.DELETE(route, handler)
			req, _ := http.NewRequest("DELETE", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
