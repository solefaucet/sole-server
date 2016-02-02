package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestSignupPayloadValidate(t *testing.T) {
	testdata := []struct {
		p       signupPayload
		canPass bool
	}{
		{
			signupPayload{BitcoinAddress: "bitcoin"},
			true,
		},
		{
			signupPayload{BitcoinAddress: "wow"},
			false,
		},
	}

	helper := func(canPass bool) (string, string) {
		if canPass {
			return "valid", "invalid"
		}
		return "invalid", "valid"
	}

	for _, v := range testdata {
		if (v.p.validate() == nil) != v.canPass {
			expected, but := helper(v.canPass)
			t.Errorf("%v should be %v but test get %v", v.p, expected, but)
		}
	}
}

func TestSignup(t *testing.T) {
	testdata := []struct {
		requestData string
		code        int
	}{
		{
			"huhu",
			400,
		},
		{
			`{"email": "", "bitcoin_address": ""}`,
			400,
		},
		{
			`{"email": "valid@email.com", "bitcoin_address": ""}`,
			400,
		},
		{
			`{"email": "valid@email.com", "bitcoin_address": "invalid"}`,
			400,
		},
		{
			`{"email": "", "bitcoin_address": "bitcoin"}`,
			400,
		},
		{
			`{"email": "valid@email.com", "bitcoin_address": "bitcoin"}`,
			200,
		},
	}

	r := gin.New()
	route := "/users"
	r.POST(route, Signup())
	for _, v := range testdata {
		req, _ := http.NewRequest("POST", route, bytes.NewBufferString(v.requestData))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != v.code {
			t.Errorf("request with %s should get status code %v but get %v", v.requestData, v.code, resp.Code)
		}
	}
}
