package v1

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterCaptcha(t *testing.T) {
	Convey("Given errored register captcha handler", t, func() {
		getCaptchaID := func() string { return "" }
		registerCaptcha := func() (string, error) { return "", fmt.Errorf("") }
		handler := RegisterCaptcha(registerCaptcha, getCaptchaID)

		Convey("When request register captcha", func() {
			route := "/captcha"
			_, resp, r := gin.CreateTestContext()
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response status code should be 500", func() {
				So(resp.Code, ShouldEqual, 500)
			})
		})
	})

	Convey("Given error free register captcha handler", t, func() {
		getCaptchaID := func() string { return "" }
		registerCaptcha := func() (string, error) { return "", nil }
		handler := RegisterCaptcha(registerCaptcha, getCaptchaID)

		Convey("When request register captcha", func() {
			route := "/captcha"
			_, resp, r := gin.CreateTestContext()
			r.POST(route, handler)
			req, _ := http.NewRequest("POST", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response status code should be 200", func() {
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}
