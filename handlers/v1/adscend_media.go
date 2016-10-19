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

type adscendMediaPayload struct {
	TransactionID string  `form:"tx_id" binding:"required"`
	UserID        int64   `form:"user_id" binding:"required,gt=0"`
	OfferID       string  `form:"offer_id"`
	Amount        float64 `form:"amount" binding:"required"`
}

// AdscendMediaCallback handles callback from adscendMedia
func AdscendMediaCallback(
	getUserByID dependencyGetUserByID,
	getAdscendMediaOffer dependencyGetAdscendMediaOffer,
	chargebackIncome dependencyChargebackIncome,
	getSystemConfig dependencyGetSystemConfig,
	createAdscendMediaIncome dependencyCreateAdscendMediaIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := adscendMediaPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":          models.EventAdscendMediaCallback,
			"query":          c.Request.URL.Query().Encode(),
			"user_id":        payload.UserID,
			"amount":         payload.Amount,
			"transaction_id": payload.TransactionID,
			"offer_id":       payload.OfferID,
		}).Debug("get adscend media callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		offer, err := getAdscendMediaOffer(payload.TransactionID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// chargeback
		if payload.Amount < 0 {
			if err := chargebackIncome(offer.IncomeID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			c.Status(http.StatusOK)
			return
		}

		// already added
		if offer != nil {
			c.String(http.StatusOK, "1")
			return
		}

		// create income adscendMedia
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeAdscendMedia,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createAdscendMediaIncome(income, payload.TransactionID, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "adscend media", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}
