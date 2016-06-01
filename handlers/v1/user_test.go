package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/solefaucet/solebtc/errors"
	"github.com/solefaucet/solebtc/models"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	validEmail   = "valid@email.cc"
	invalidEmail = "invalid@.ee.cc"
)

func TestSignup(t *testing.T) {
	requestDataJSON := func(email string) []byte {
		raw, _ := json.Marshal(map[string]interface{}{
			"email":      email,
			"address":    "address",
			"referer_id": 2,
		})
		return raw
	}

	testdata := []struct {
		when            string
		requestData     []byte
		code            int
		getUserByID     dependencyGetUserByID
		createUser      dependencyCreateUser
		validateAddress func(string) (bool, error)
	}{
		{
			"invalid json data",
			[]byte("huhu"),
			400,
			nil,
			nil,
			nil,
		},
		{
			"invalid email",
			requestDataJSON(invalidEmail),
			400,
			nil,
			nil,
			nil,
		},
		{
			"invalid address",
			requestDataJSON(validEmail),
			400,
			nil,
			nil,
			func(string) (bool, error) { return false, nil },
		},
		{
			"duplicate email",
			requestDataJSON(validEmail),
			409,
			mockGetUserByID(models.User{}, nil),
			mockCreateUser(errors.ErrDuplicatedEmail),
			func(string) (bool, error) { return true, nil },
		},
		{
			"valid email, but create user unknown error",
			requestDataJSON(validEmail),
			500,
			mockGetUserByID(models.User{}, nil),
			mockCreateUser(fmt.Errorf("")),
			func(string) (bool, error) { return true, nil },
		},
		{
			"valid email",
			requestDataJSON(validEmail),
			200,
			mockGetUserByID(models.User{}, nil),
			mockCreateUser(nil),
			func(string) (bool, error) { return true, nil },
		},
	}

	for _, v := range testdata {
		Convey("Given Signup controller", t, func() {
			handler := Signup(v.validateAddress, v.createUser, v.getUserByID)

			Convey(fmt.Sprintf("When request with %s", v.when), func() {
				route := "/users"
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

func TestVerifyEmail(t *testing.T) {
	Convey("Given verify email controller with expired session and errored getSessionByToken dependency", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{}, fmt.Errorf(""))
		handler := VerifyEmail(getSessionByToken, nil, nil)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given verify email controller with expired session and getSessionByToken dependency", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{}, nil)
		handler := VerifyEmail(getSessionByToken, nil, nil)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 401", func() {
				So(resp.Code, ShouldEqual, 401)
			})
		})
	})

	Convey("Given verify email controller with errored getUserByID dependency", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{UpdatedAt: time.Now()}, nil)
		getUserByID := mockGetUserByID(models.User{}, fmt.Errorf(""))
		handler := VerifyEmail(getSessionByToken, getUserByID, nil)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given verify email controller with banned user status", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{UpdatedAt: time.Now()}, nil)
		getUserByID := mockGetUserByID(models.User{Status: models.UserStatusBanned}, nil)
		handler := VerifyEmail(getSessionByToken, getUserByID, nil)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 403", func() {
				So(resp.Code, ShouldEqual, 403)
			})
		})
	})

	Convey("Given verify email controller with errored updateUser dependency", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{UpdatedAt: time.Now()}, nil)
		getUserByID := mockGetUserByID(models.User{}, nil)
		updateUserStatus := mockUpdateUserStatus(fmt.Errorf(""))
		handler := VerifyEmail(getSessionByToken, getUserByID, updateUserStatus)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given verify email controller with correct dependencies injected", t, func() {
		getSessionByToken := mockGetSessionByToken(models.Session{UpdatedAt: time.Now()}, nil)
		getUserByID := mockGetUserByID(models.User{}, nil)
		updateUserStatus := mockUpdateUserStatus(nil)
		handler := VerifyEmail(getSessionByToken, getUserByID, updateUserStatus)

		Convey("When verify email", func() {
			route := "/users/1/status"
			_, resp, r := gin.CreateTestContext()
			r.PUT(route, handler)
			req, _ := http.NewRequest("PUT", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestGetUserInfo(t *testing.T) {
	Convey("Given get user info controller with errored getUserByID dependency", t, func() {
		getUserByID := mockGetUserByID(models.User{}, errors.ErrNotFound)
		handler := UserInfo(getUserByID)

		Convey("When get user info", func() {
			route := "/users"
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

	Convey("Given get user info controller with correctly dependencies injected", t, func() {
		getUserByID := mockGetUserByID(models.User{}, nil)
		handler := UserInfo(getUserByID)

		Convey("When get user info", func() {
			route := "/users"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestGetReferees(t *testing.T) {
	Convey("Given referee list handler with errored dependency", t, func() {
		handler := RefereeList(mockGetReferees(nil, fmt.Errorf("")), nil)

		Convey("When get reward list with invalid offset", func() {
			route := "/users/referees"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route+"?offset=3i", nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})

		Convey("When get referee list", func() {
			route := "/users/referees"
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

	Convey("Given referee list controller with correct dependencies injected", t, func() {
		handler := RefereeList(mockGetReferees(nil, nil), func(int64) (int64, error) { return 0, nil })

		Convey("When get referee list", func() {
			route := "/users/referees"
			_, resp, r := gin.CreateTestContext()
			r.Use(func(c *gin.Context) {
				c.Set("auth_token", models.AuthToken{})
			})
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
