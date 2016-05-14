package middlewares

import (
	"net/http"

	"github.com/freeusd/solebtc/errors"
	"github.com/gin-gonic/gin"
)

type dependencyValidateCaptcha func(challenge, validate, seccode string) (bool, error)

// CaptchaValidationRequired middleware validates captcha
func CaptchaValidationRequired(validateCaptcha dependencyValidateCaptcha) gin.HandlerFunc {
	return func(c *gin.Context) {
		challenge := c.Request.Header.Get("X-Geetest-Challenge")
		validate := c.Request.Header.Get("X-Geetest-Validate")
		seccode := c.Request.Header.Get("X-Geetest-Seccode")

		valid, err := validateCaptcha(challenge, validate, seccode)
		if err != nil {
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
