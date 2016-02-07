package v1

import (
	"net/http"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/satori/go.uuid"
	"github.com/freeusd/solebtc/constant"
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

type (
	requestVerifyEmailDependencyGetUserByID   func(int) (models.User, *errors.Error)
	requestVerifyEmailDependencyUpsertSession func(models.Session) *errors.Error
	requestVerifyEmailDependencySendEmail     func(recipients []string, subject string, html string) *errors.Error
)

// RequestVerifyEmail send verification url to user via email
func RequestVerifyEmail(
	getUserByID requestVerifyEmailDependencyGetUserByID,
	upsertSession requestVerifyEmailDependencyUpsertSession,
	sendEmail requestVerifyEmailDependencySendEmail,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// get user
		user, err := getUserByID(authToken.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// upsert session
		if err := upsertSession(models.Session{
			UserID: authToken.UserID,
			Token:  uuid.NewV4().String(),
			Type:   constant.SessionTypeVerifyEmail,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// send email
		// FIXME: should be fix later on
		subject := "Verify your email in SoleBTC"
		html := ""
		if err := sendEmail([]string{user.Email}, subject, html); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
