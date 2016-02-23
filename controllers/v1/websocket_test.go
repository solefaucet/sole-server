package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gorilla/websocket"
	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
)

func TestWebsocket(t *testing.T) {
	Convey("Given a webserver with websocket controller", t, func() {
		handler := Websocket(func(*websocket.Conn) {})

		Convey("Websocket connect", func() {
			route := "/websocket"
			_, _, r := gin.CreateTestContext()
			r.GET(route, handler)
			server := httptest.NewServer(r)

			_, resp, _ := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+route, nil)

			Convey("Response code should be 101", func() {
				So(resp.StatusCode, ShouldEqual, 101)
			})

			Reset(func() {
				server.Close()
			})
		})
	})

	Convey("Given websocket controller", t, func() {
		handler := Websocket(func(*websocket.Conn) {})

		Convey("HTTP connect", func() {
			route := "/websocket"
			_, resp, r := gin.CreateTestContext()
			r.GET(route, handler)
			req, _ := http.NewRequest("GET", route, nil)
			r.ServeHTTP(resp, req)

			Convey("Response code should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})
}
