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

type personalyPayload struct {
	UserID  int64   `form:"user_id" binding:"required,gt=0"`
	OfferID string  `form:"offer_id" binding:"required"`
	Amount  float64 `form:"amount" binding:"required,gt=0"`
}

// PersonalyCallback handles callback from personaly
func PersonalyCallback(
	getUserByID dependencyGetUserByID,
	getNumberOfPersonalyOffers dependencyGetNumberOfPersonalyOffers,
	getSystemConfig dependencyGetSystemConfig,
	createPersonalyIncome dependencyCreatePersonalyIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := personalyPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":    models.EventPersonalyCallback,
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

		count, err := getNumberOfPersonalyOffers(payload.OfferID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.String(http.StatusOK, "1")
			return
		}

		// create income personaly
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypePersonaly,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createPersonalyIncome(income, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "personaly", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}
