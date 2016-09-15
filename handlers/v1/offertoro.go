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

type offertoroPayload struct {
	TransactionID string  `form:"id" binding:"required"`
	UserID        int64   `form:"user_id" binding:"required,gt=0"`
	OfferID       string  `form:"oid"`
	Amount        float64 `form:"amount" binding:"required,gt=0"`
}

// OffertoroCallback handles callback from offertoro
func OffertoroCallback(
	getUserByID dependencyGetUserByID,
	getNumberOfOffertoroOffers dependencyGetNumberOfOffertoroOffers,
	getSystemConfig dependencyGetSystemConfig,
	createOffertoroIncome dependencyCreateOffertoroIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := offertoroPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":          models.EventOffertoroCallback,
			"query":          c.Request.URL.Query().Encode(),
			"user_id":        payload.UserID,
			"amount":         payload.Amount,
			"transaction_id": payload.TransactionID,
			"offer_id":       payload.OfferID,
		}).Debug("get offertoro callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfOffertoroOffers(payload.TransactionID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.String(http.StatusOK, "1")
			return
		}

		// create income offertoro
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeOffertoro,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createOffertoroIncome(income, payload.TransactionID, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "offertoro", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}
