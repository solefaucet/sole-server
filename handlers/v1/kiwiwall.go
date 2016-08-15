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

type kiwiwallPayload struct {
	TransactionID string  `form:"trans_id" binding:"required"`
	UserID        int64   `form:"sub_id" binding:"required,gt=0"`
	Amount        float64 `form:"amount" binding:"required,gt=0"`
	OfferID       string  `form:"offer_id"`
	OfferName     string  `form:"offer_name"`
}

// KiwiwallCallback handles callback from kiwiwall
func KiwiwallCallback(
	getUserByID dependencyGetUserByID,
	getNumberOfKiwiwallOffers dependencyGetNumberOfKiwiwallOffers,
	getSystemConfig dependencyGetSystemConfig,
	createKiwiwallIncome dependencyCreateKiwiwallIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := kiwiwallPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":          models.EventKiwiwallCallback,
			"query":          c.Request.URL.Query().Encode(),
			"user_id":        payload.UserID,
			"amount":         payload.Amount,
			"transaction_id": payload.TransactionID,
			"offer_id":       payload.OfferID,
			"offer_name":     payload.OfferName,
		}).Debug("get kiwiwall callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfKiwiwallOffers(payload.TransactionID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.String(http.StatusOK, "1")
			return
		}

		// create income kiwiwall
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeKiwiwall,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createKiwiwallIncome(income, payload.TransactionID, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "kiwiwall", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}
