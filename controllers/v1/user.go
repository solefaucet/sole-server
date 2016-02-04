package v1

import (
	"fmt"
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/utils"
)

type signupPayload struct {
	Email          string `json:"email" binding:"required,email"`
	BitcoinAddress string `json:"bitcoin_address" binding:"required"`
}

func (p *signupPayload) validate() error {
	if ok, err := utils.ValidateBitcoinAddress(p.BitcoinAddress); err != nil || !ok {
		e := errors.Error{
			ErrCode:   errors.ErrCodeInvalidBitcoinAddress,
			ErrString: fmt.Sprintf("%s is invalid bitcoin address", p.BitcoinAddress),
		}
		if err != nil {
			e.ErrStringForLogging = fmt.Sprintf("validate bitcoin address error: %v", err)
		}
		return e
	}

	return nil
}

// Signup creates a new user with unique email, bitcoin address
func Signup(dependencies ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := signupPayload{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}
		if err := payload.validate(); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Status(200)
	}
}
