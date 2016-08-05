package v1

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/models"
)

// RegisterCaptcha register get challenge from geetest
func RegisterCaptcha(
	registerCaptcha dependencyRegisterCaptcha,
	getCaptchaID dependencyGetCaptchaID,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithField("event", models.EventRegisterCaptcha).Info("registering geetest captcha")
		captchaID := getCaptchaID()
		challenge, err := registerCaptcha()
		if err != nil {
			c.AbortWithError(http.StatusServiceUnavailable, err)
			return
		}

		c.JSON(http.StatusOK, map[string]string{
			"captcha_id": captchaID,
			"challenge":  challenge,
		})
	}
}
