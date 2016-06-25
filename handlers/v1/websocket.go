package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/solefaucet/sole-server/models"
)

var defaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow any origin to connect
}

// Websocket entry
func Websocket(
	getUsersOnline dependencyGetUsersOnline,
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

		// increment usersOnline by 1
		usersOnline := getUsersOnline() + 1

		// send msg to client
		conn.WriteJSON(models.WebsocketMessage{
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
