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

const (
	validBTCAddr   = "1EFJFaeATfp2442TGcHS5mgadXJjsSSP2T"
	invalidBTCAddr = "bitcoin"

	validEmail   = "valid@email.cc"
	invalidEmail = "invalid@.ee.cc"
)

type mockCreateUserService struct {
	createUserService
	err errors.Error
}

func (m mockCreateUserService) CreateUser(models.User) errors.Error {
	return m.err
}

func newMockCreateUserService(err errors.Error) mockCreateUserService {
	return mockCreateUserService{err: err}
}

func TestSignup(t *testing.T) {
	requestDataJSON := func(email, btcAddr string) []byte {
		raw, _ := json.Marshal(map[string]interface{}{
			"email":           email,
			"bitcoin_address": btcAddr,
		})
		return raw
	}

	testdata := []struct {
		when        string
		requestData []byte
		code        int
		mcus        createUserService
	}{
		{
			"invalid json data",
			[]byte("huhu"),
			400,
			nil,
		},
		{
			"invalid email, invalid bitcoin address",
			requestDataJSON(invalidEmail, invalidBTCAddr),
			400,
			nil,
		},
		{
			"valid email, invalid bitcoin address",
			requestDataJSON(validEmail, invalidEmail),
			400,
			nil,
		},
		{
			"invalid email, valid bitcoin address",
			requestDataJSON(invalidEmail, validBTCAddr),
			400,
			nil,
		},
		{
			"duplicate email, valid bitcoin address",
			requestDataJSON(validEmail, validBTCAddr),
			409,
			newMockCreateUserService(errors.Error{
				ErrCode: errors.ErrCodeDuplicateEmail,
			}),
		},
		{
			"valid email, duplicate bitcoin address",
			requestDataJSON(validEmail, validBTCAddr),
			409,
			newMockCreateUserService(errors.Error{
				ErrCode: errors.ErrCodeDuplicateBitcoinAddress,
			}),
		},
		{
			"valid email, valid bitcoin address, but unknow error",
			requestDataJSON(validEmail, validBTCAddr),
			500,
			newMockCreateUserService(errors.Error{
				ErrCode: errors.ErrCodeUnknown,
			}),
		},
		{
			"valid email, valid bitcoin address",
			requestDataJSON(validEmail, validBTCAddr),
			200,
			newMockCreateUserService(errors.Nil),
		},
	}

	for _, v := range testdata {
		Convey("Given Signup controller", t, func() {
			s := Signup(v.mcus)
			route := "/users"

			Convey(fmt.Sprintf("When request with %s", v.when), func() {
				_, resp, r := gin.CreateTestContext()
				r.POST(route, s)
				req, _ := http.NewRequest("POST", route, bytes.NewBuffer(v.requestData))
				r.ServeHTTP(resp, req)

				Convey("Response code should equal", func() {
					So(resp.Code, ShouldEqual, v.code)
				})
			})
		})
	}
}
