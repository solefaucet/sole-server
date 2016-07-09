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

type clixwallPayload struct {
	Password  string  `form:"pwd" binding:"required"`
	UserID    int64   `form:"u" binding:"required,gt=0"`
	Amount    float64 `form:"c" binding:"required,gt=0"`
	OfferID   string  `form:"cid" binding:"required"`
	OfferName string  `form:"cname"`
	Status    int64   `form:"s" binding:"required,eq=1"` // Must be status(1) -> 1. credited 2. debit
	Type      int64   `form:"t" binding:"required,eq=2"` // Must be points(2) -> 1. cash 2. points
}

// ClixwallCallback handles callback from clixwall
func ClixwallCallback(
	secretPassword string,
	getUserByID dependencyGetUserByID,
	getNumberOfClixwallOffers dependencyGetNumberOfClixwallOffers,
	getSystemConfig dependencyGetSystemConfig,
	createClixwallIncome dependencyCreateClixwallIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := clixwallPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		if secretPassword == "" || secretPassword != payload.Password {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":      models.EventClixwallCallback,
			"user_id":    payload.UserID,
			"amount":     payload.Amount,
			"offer_id":   payload.OfferID,
			"offer_name": payload.OfferName,
		}).Debug("get clixwall callback")

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfClixwallOffers(payload.OfferID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.Status(http.StatusOK)
			return
		}

		// create income clixwall
		amount := payload.Amount
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeClixwall,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createClixwallIncome(income, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, payload.Amount, "clixwall", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.Status(http.StatusOK)
	}
}
