package v1

import (
	"fmt"
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/utils"
)

type signupPayload struct {
	Email          string `json:"email" binding:"required,email"`
	BitcoinAddress string `json:"bitcoin_address" binding:"required"`
}

func (p *signupPayload) validate() error {
	if ok, err := utils.ValidateBitcoinAddress(p.BitcoinAddress); err != nil || !ok {
		e := errors.New(errors.ErrCodeInvalidBitcoinAddress)
		e.ErrString = fmt.Sprintf("%s is invalid bitcoin address", p.BitcoinAddress)
		if err != nil {
			e.ErrStringForLogging = fmt.Sprintf("validate bitcoin address error: %v", err)
		}
		return e
	}

	return nil
}

type signupDependencyCreateUser func(models.User) *errors.Error

func userWithSignupPayload(p signupPayload) models.User {
	return models.User{
		Email:          p.Email,
		BitcoinAddress: p.BitcoinAddress,
	}
}

// Signup creates a new user with unique email, bitcoin address
func Signup(createUser signupDependencyCreateUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := signupPayload{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}
		if err := payload.validate(); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user := userWithSignupPayload(payload)
		if err := createUser(user); err != nil {
			switch err.ErrCode {
			case errors.ErrCodeDuplicateEmail:
				err.ErrString = fmt.Sprintf("Email %s is duplicated", payload.Email)
				c.AbortWithError(http.StatusConflict, err)
			case errors.ErrCodeDuplicateBitcoinAddress:
				err.ErrString = fmt.Sprintf("Bitcoin address %s is duplicated", payload.BitcoinAddress)
				c.AbortWithError(http.StatusConflict, err)
			default:
				err.ErrString = "Internal server error"
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
