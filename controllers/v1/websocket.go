package v1

import (
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gorilla/websocket"
)

var defaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow any origin to connect
}

// Websocket entry
func Websocket(putConn dependencyPutConn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade http conn to ws conn
		conn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		putConn(conn)
	}
}
