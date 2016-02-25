package v1

import (
	"encoding/json"
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/freeusd/solebtc/models"
)

var defaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow any origin to connect
}

// Websocket entry
func Websocket(
	getUsersOnline dependencyGetUsersOnline,
	getConfig dependencyGetSystemConfig,
	getLatestIncomes dependencyGetLatestIncomes,
	broadcast dependencyBroadcast,
	putConn dependencyPutConn,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade http conn to ws conn
		conn, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// DO NOT allow client side send msg other than PING
		conn.SetReadLimit(1)

		// increment usersOnline by 1
		usersOnline := getUsersOnline() + 1

		// send msg to client
		conn.WriteJSON(models.WebsocketMessage{
			BitcoinPrice:  getConfig().BitcoinPrice,
			UsersOnline:   usersOnline,
			LatestIncomes: getLatestIncomes(),
		})

		// broadcast users_online to all other clients
		broadcastMsg, _ := json.Marshal(models.WebsocketMessage{
			UsersOnline: usersOnline,
		})
		broadcast(broadcastMsg)

		// put connection in hub
		putConn(conn)
	}
}
