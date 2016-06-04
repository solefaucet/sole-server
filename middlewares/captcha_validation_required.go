package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

type dependencyValidateCaptcha func(challenge, validate, seccode string) (bool, error)

// CaptchaValidationRequired middleware validates captcha
func CaptchaValidationRequired(validateCaptcha dependencyValidateCaptcha) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithField("event", models.EventValidateCaptcha).Info("validating geetest captcha")
		challenge := c.Request.Header.Get("X-Geetest-Challenge")
		validate := c.Request.Header.Get("X-Geetest-Validate")
		seccode := c.Request.Header.Get("X-Geetest-Seccode")

		valid, err := validateCaptcha(challenge, validate, seccode)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"event":     models.EventValidateCaptcha,
				"challenge": challenge,
				"validate":  validate,
				"seccode":   seccode,
			}).Warn(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if !valid {
			c.AbortWithError(http.StatusBadRequest, errors.ErrInvalidCaptcha)
			return
		}

		c.Next()
	}
}
