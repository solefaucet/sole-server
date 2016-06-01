package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/solefaucet/solebtc/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWebsocket(t *testing.T) {
	Convey("Given a webserver with websocket controller", t, func() {
		handler := Websocket(
			mockGetUsersOnline(1),
			mockGetLatestIncomes([]interface{}{123.456, "hello"}),
			mockBroadcast(),
			mockPutConn(),
		)

		Convey("When websocket connect", func() {
			route := "/websocket"
			_, _, r := gin.CreateTestContext()
			r.GET(route, handler)
			server := httptest.NewServer(r)

			conn, resp, _ := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+route, nil)

			Convey("Response code should be 101", func() {
				So(resp.StatusCode, ShouldEqual, 101)
			})

			m := models.WebsocketMessage{}
			conn.ReadJSON(&m)
			Convey("Receive data through websocket should resemble", func() {
				So(m, ShouldResemble, models.WebsocketMessage{
					UsersOnline:   2,
					LatestIncomes: []interface{}{123.456, "hello"},
				})
			})

			Reset(func() {
				server.Close()
				conn.Close()
			})
		})
	})

	Convey("Given websocket controller", t, func() {
		handler := Websocket(nil, nil, nil, nil)

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
