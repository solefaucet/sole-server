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

func mockLoginDependencyGetUserByEmail(err *errors.Error) loginDependencyGetUserByEmail {
	return func(string) (models.User, *errors.Error) {
		return models.User{}, err
	}
}

func mockLoginDependencyCreateAuthToken(err *errors.Error) loginDependencyCreateAuthToken {
	return func(models.AuthToken) *errors.Error {
		return err
	}
}

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
		getUserByEmail  loginDependencyGetUserByEmail
		createAuthToken loginDependencyCreateAuthToken
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
			"non existing email",
			requestDataJSON(validEmail),
			404,
			mockLoginDependencyGetUserByEmail(errors.New(errors.ErrCodeNotFound)),
			nil,
		},
		{
			"valid email, unknown error",
			requestDataJSON(validEmail),
			500,
			mockLoginDependencyGetUserByEmail(errors.New(errors.ErrCodeUnknown)),
			nil,
		},
		{
			"valid existing email, but unknown error",
			requestDataJSON(validEmail),
			500,
			mockLoginDependencyGetUserByEmail(nil),
			mockLoginDependencyCreateAuthToken(errors.New(errors.ErrCodeUnknown)),
		},
		{
			"valid existing email",
			requestDataJSON(validEmail),
			201,
			mockLoginDependencyGetUserByEmail(nil),
			mockLoginDependencyCreateAuthToken(nil),
		},
	}

	for _, v := range testdata {
		Convey("Given Login controller", t, func() {
			s := Login(v.getUserByEmail, v.createAuthToken)

			Convey(fmt.Sprintf("When request with %s", v.when), func() {
				route := "/auth_tokens"
				_, resp, r := gin.CreateTestContext()
				r.POST(route, s)
				req, _ := http.NewRequest("POST", route, bytes.NewBuffer(v.requestData))
				r.ServeHTTP(resp, req)

				Convey(fmt.Sprintf("Response code should be equal to %d", v.code), func() {
					So(resp.Code, ShouldEqual, v.code)
				})
			})
		})
	}
}
