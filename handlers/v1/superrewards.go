package v1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/sole-server/models"
)

type superrewardPayload struct {
	TransactionID string  `form:"id" binding:"required"`
	UserID        int64   `form:"uid" binding:"required,gt=0"`
	OfferID       string  `form:"oid"`
	Amount        float64 `form:"new" binding:"required,gt=0"`
	Total         float64 `form:"total"`
	Signature     string  `form:"sig" binding:"required"`
}

// SuperrewardsCallback handles callback from superrewards
func SuperrewardsCallback(
	secretKey string,
	getUserByID dependencyGetUserByID,
	getNumberOfSuperrewardsOffers dependencyGetNumberOfSuperrewardsOffers,
	getSystemConfig dependencyGetSystemConfig,
	createSuperrewardsIncome dependencyCreateSuperrewardsIncome,
	broadcast dependencyBroadcast,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := superrewardPayload{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":           models.EventSuperrewardsCallback,
			"query":           c.Request.URL.Query().Encode(),
			"user_id":         payload.UserID,
			"amount":          payload.Amount,
			"transaction_id":  payload.TransactionID,
			"offer_id":        payload.OfferID,
			"user_total_earn": payload.Total,
		}).Debug("get superrewards callback")

		if !validateSuperrewardsRequest(payload, secretKey) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		user, err := getUserByID(payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfSuperrewardsOffers(payload.TransactionID, payload.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if count > 0 {
			c.String(http.StatusOK, "1")
			return
		}

		// create income superrewards
		amount := payload.Amount / 1e8
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeSuperrewards,
			Income:        amount,
			RefererIncome: amount * getSystemConfig().RefererRewardRate,
		}
		if err := createSuperrewardsIncome(income, payload.TransactionID, payload.OfferID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// broadcast delta income to all clients
		deltaIncome := struct {
			Address string    `json:"address"`
			Amount  float64   `json:"amount"`
			Type    string    `json:"type"`
			Time    time.Time `json:"time"`
		}{user.Address, amount, "superrewards", time.Now()}
		msg, _ := json.Marshal(models.WebsocketMessage{DeltaIncome: deltaIncome})
		broadcast(msg)

		c.String(http.StatusOK, "1")
	}
}

func validateSuperrewardsRequest(payload superrewardPayload, secretKey string) bool {
	data := fmt.Sprintf("%v:%v:%v:%v", payload.TransactionID, payload.Amount, payload.UserID, secretKey)
	if sign := fmt.Sprintf("%x", md5.Sum([]byte(data))); sign != payload.Signature {
		logrus.WithFields(logrus.Fields{
			"event":   models.EventSuperrewardsInvalidSignature,
			"payload": payload,
			"sign":    sign,
		}).Debug("signature not match")
		return false
	}
	return true
}
