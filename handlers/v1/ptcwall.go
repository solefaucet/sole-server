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
	UserID   int64   `form:"usr" binding:"required,gt=0"`
	Credited int64   `form:"c" binding:"required,eq=1"` // Must be credited(1) -> 1. credited 2. reversed
	Type     int64   `form:"t" binding:"required,eq=2"` // Must be points(2) -> 1. cash 2. points
	Amount   float64 `form:"r" binding:"required,gt=0"`
}

// PtcwallCallback handles callback from ptcwall
func PtcwallCallback(
	getUserByID dependencyGetUserByID,
	getSystemConfig dependencyGetSystemConfig,
	createPtcwallIncome dependencyCreatePtcwallIncome,
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
			"amount":  payload.Amount,
		}).Debug("get ptcwall callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// create income ptcwall
		amount := payload.Amount
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypePtcwall,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createPtcwallIncome(income); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, payload.Amount, "ptcwall", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.Status(http.StatusOK)
	}
}
