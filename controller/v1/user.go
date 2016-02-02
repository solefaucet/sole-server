package v1

import (
	"errors"
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

type signupPayload struct {
	Email          string `json:"email" binding:"required,email"`
	BitcoinAddress string `json:"bitcoin_address" binding:"required"`
}

func (p *signupPayload) validate() error {
	if p.BitcoinAddress != "bitcoin" {
		return errors.New("invalid bitcoin address")
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
