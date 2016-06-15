package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

/*
immediate:

0：非即时返利活动,处于待审核状态；
1：即时返利活动，需发放奖励给会员；
2：非即时返利活动，审核通过，重新回传，发放奖励给会员；
3：非即时返利活动，审核不通过，重新回传，不发放奖励；

errno:

offerwow-01: 出现空参数
offerwow-02: 网站id不存在
offerwow-03: uid会员不存在
offerwow-04: 已发放奖励的Eventid重复
offerwow-05: immediate=0
offerwow-06: immediate=3
*/

// offerwow const
const (
	offerwowWebsiteID     = "2626"
	offerwowImmediate0    = "0"
	offerwowImmediate1    = "1"
	offerwowImmediate2    = "2"
	offerwowImmediate3    = "3"
	offerwowStatusSuccess = "success"
	offerwowStatusFailure = "failure"
	offerwowErrno01       = "offerwow-01"
	offerwowErrno02       = "offerwow-02"
	offerwowErrno03       = "offerwow-03"
	offerwowErrno04       = "offerwow-04"
	offerwowErrno05       = "offerwow-05"
	offerwowErrno06       = "offerwow-06"
)

type offerwowPayload struct {
	UserID      int64   `form:"memberid"`
	Amount      float64 `form:"point"`
	EventID     string  `form:"eventid"`
	WebsiteID   string  `form:"websiteid"`
	Immediate   string  `form:"immediate"`
	ProgramName string  `form:"programname"`
	Sign        string  `form:"sign"` // NOTE: should be required in production
}

type offerwowResponse struct {
	UserID    int64
	Amount    float64
	EventID   string
	Immediate string
	Status    string
	Error     string
}

func (r offerwowResponse) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"memberid":  fmt.Sprint(r.UserID),
		"point":     fmt.Sprint(r.Amount),
		"websiteid": offerwowWebsiteID,
		"eventid":   r.EventID,
		"Immediate": r.Immediate,
		"status":    r.Status,
	}
	if r.Error != "" {
		m["errno"] = r.Error
	}
	return json.Marshal(m)
}

// OfferwowCallback handles callback from offerwow
func OfferwowCallback(
	getUserByID dependencyGetUserByID,
	getOfferwowEventByID dependencyGetOfferwowEventByID,
	getSystemConfig dependencyGetSystemConfig,
	createOfferwowIncome dependencyCreateOfferwowIncome,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := offerwowPayload{}
		c.BindWith(&payload, binding.Form)

		responseAndLog := func(status, err string) {
			response := offerwowResponse{
				UserID:    payload.UserID,
				Amount:    payload.Amount,
				EventID:   payload.EventID,
				Immediate: payload.Immediate,
				Status:    status,
				Error:     err,
			}
			c.JSON(http.StatusOK, response)
			logrus.WithFields(logrus.Fields{
				"event":    models.EventOfferwowCallback,
				"query":    c.Request.URL.Query().Encode(),
				"payload":  payload,
				"response": response,
			}).Debug("get offerwow callback")
		}

		if payload.UserID == 0 || payload.Amount == 0 || payload.EventID == "" || payload.WebsiteID == "" || payload.Immediate == "" {
			responseAndLog(offerwowStatusFailure, offerwowErrno01)
			return
		}

		if payload.WebsiteID != offerwowWebsiteID {
			responseAndLog(offerwowStatusFailure, offerwowErrno02)
			return
		}

		// check if user exists
		user, err := getUserByID(payload.UserID)
		if err != nil {
			switch err {
			case errors.ErrNotFound:
				responseAndLog(offerwowStatusFailure, offerwowErrno03)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		// check if eventid duplicates
		if _, err := getOfferwowEventByID(payload.EventID); err != nil {
			switch err {
			case errors.ErrNotFound:
				responseAndLog(offerwowStatusFailure, offerwowErrno04)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		if payload.Immediate == offerwowImmediate0 {
			responseAndLog(offerwowStatusFailure, offerwowErrno05)
			return
		}

		if payload.Immediate == offerwowImmediate3 {
			responseAndLog(offerwowStatusFailure, offerwowErrno06)
			return
		}

		// create income offerwow
		income := models.Income{
			UserID:        user.ID,
			RefererID:     user.RefererID,
			Type:          models.IncomeTypeOfferwow,
			Income:        payload.Amount / 1e8,
			RefererIncome: getSystemConfig().RefererRewardRate * payload.Amount,
		}
		if err := createOfferwowIncome(income, payload.EventID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		responseAndLog(offerwowStatusSuccess, "")
	}
}
