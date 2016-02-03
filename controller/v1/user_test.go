package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	validBTCAddr   = "1EFJFaeATfp2442TGcHS5mgadXJjsSSP2T"
	invalidBTCAddr = "bitcoin"

	validEmail   = "valid@email.cc"
	invalidEmail = "invalid@.ee.cc"
)

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
	}{
		{
			"invalid json data",
			[]byte("huhu"),
			400,
		},
		{
			"invalid email, invalid bitcoin address",
			requestDataJSON(invalidEmail, invalidBTCAddr),
			400,
		},
		{
			"valid email, invalid bitcoin address",
			requestDataJSON(validEmail, invalidEmail),
			400,
		},
		{
			"invalid email, valid bitcoin address",
			requestDataJSON(invalidEmail, validBTCAddr),
			400,
		},
		{
			"valid email, valid bitcoin address",
			requestDataJSON(validEmail, validBTCAddr),
			200,
		},
	}

	Convey("Given Signup controller", t, func() {
		s := Signup()
		route := "/users"

		for _, v := range testdata {
			Convey(fmt.Sprintf("When request with %s", v.when), func() {
				_, resp, r := gin.CreateTestContext()
				r.POST(route, s)
				req, _ := http.NewRequest("POST", route, bytes.NewBuffer(v.requestData))
				r.ServeHTTP(resp, req)

				Convey("Response code should equal", func() {
					So(resp.Code, ShouldEqual, v.code)
				})
			})
		}
	})
}
