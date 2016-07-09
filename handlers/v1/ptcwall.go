package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/sole-server/models"
)

type ptcwallPayload struct {
	UserID      int64   `form:"usr" binding:"required,gt=0"`
	Credited    int64   `form:"c" binding:"required,eq=1"` // Must be credited(1) -> 1. credited 2. reversed
	Type        int64   `form:"t" binding:"required,eq=2"` // Must be points(2) -> 1. cash 2. points
	Rate        float64 `form:"rate" binding:"required,gt=0"`
	Transaction string  `form:"none"`
}

// PTCWallCallback handles callback from ptcwall
func PTCWallCallback(
	getUserByID dependencyGetUserByID,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := ptcwallPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":   models.EventPTCWallCallback,
			"query":   c.Request.URL.Query().Encode(),
			"user_id": payload.UserID,
		}).Debug("get ptcwall callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, payload.Rate, "ptcwall", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.Status(http.StatusOK)
	}
}
