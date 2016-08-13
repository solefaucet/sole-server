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

type trialpayPayload struct {
	UserID  int64   `form:"sid" binding:"required,gt=0"`
	OfferID string  `form:"oid" binding:"required"`
	Amount  float64 `form:"reward_amount" binding:"required,gt=0"`
}

// TrialpayCallback handles callback from trialpay
func TrialpayCallback(
	getUserByID dependencyGetUserByID,
	getNumberOfTrialpayOffers dependencyGetNumberOfTrialpayOffers,
	getSystemConfig dependencyGetSystemConfig,
	createTrialpayIncome dependencyCreateTrialpayIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := trialpayPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":    models.EventTrialpayCallback,
			"query":    c.Request.URL.Query().Encode(),
			"user_id":  payload.UserID,
			"amount":   payload.Amount,
			"offer_id": payload.OfferID,
		}).Debug("get superrewards callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfTrialpayOffers(payload.OfferID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		// create income trialpay
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeTrialpay,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createTrialpayIncome(income, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "trialpay", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.Status(http.StatusOK)
	}
}
