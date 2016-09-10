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

type adgateMediaPayload struct {
	TransactionID string  `form:"tx_id" binding:"required"`
	UserID        int64   `form:"user_id" binding:"required,gt=0"`
	OfferID       string  `form:"offer_id"`
	Amount        float64 `form:"point_value" binding:"required,gt=0"`
}

// AdgateMediaCallback handles callback from adgateMedia
func AdgateMediaCallback(
	getUserByID dependencyGetUserByID,
	getNumberOfAdgateMediaOffers dependencyGetNumberOfAdgateMediaOffers,
	getSystemConfig dependencyGetSystemConfig,
	createAdgateMediaIncome dependencyCreateAdgateMediaIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := adgateMediaPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":          models.EventAdgateMediaCallback,
			"query":          c.Request.URL.Query().Encode(),
			"user_id":        payload.UserID,
			"amount":         payload.Amount,
			"transaction_id": payload.TransactionID,
			"offer_id":       payload.OfferID,
		}).Debug("get adgate media callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfAdgateMediaOffers(payload.TransactionID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.String(http.StatusOK, "1")
			return
		}

		// create income adgateMedia
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeAdgateMedia,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createAdgateMediaIncome(income, payload.TransactionID, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "adgate media", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}
