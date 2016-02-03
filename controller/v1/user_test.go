package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

const (
	validBTCAddr   = "1EFJFaeATfp2442TGcHS5mgadXJjsSSP2T"
	invalidBTCAddr = "bitcoin"
	emptyBTCAddr   = ""

	validEmail   = "valid@email.cc"
	invalidEmail = "invalid@.ee.cc"
	emptyEmail   = ""
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
		requestData []byte
		code        int
	}{
		{
			[]byte("huhu"),
			400,
		},
		{
			requestDataJSON(emptyEmail, emptyBTCAddr),
			400,
		},
		{
			requestDataJSON(validEmail, emptyBTCAddr),
			400,
		},
		{
			requestDataJSON(validEmail, invalidBTCAddr),
			400,
		},
		{
			requestDataJSON(emptyEmail, validBTCAddr),
			400,
		},
		{
			requestDataJSON(validEmail, validBTCAddr),
			200,
		},
	}

	for _, v := range testdata {
		route := "/users"
		_, resp, r := gin.CreateTestContext()
		r.POST(route, Signup())
		req, _ := http.NewRequest("POST", route, bytes.NewBuffer(v.requestData))
		r.ServeHTTP(resp, req)

		if resp.Code != v.code {
			t.Errorf("request with %s should get status code %v but get %v", v.requestData, v.code, resp.Code)
		}
	}
}
