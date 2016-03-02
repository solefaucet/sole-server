package v1

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"

	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/satori/go.uuid"
	"github.com/freeusd/solebtc/models"
)

// RequestVerifyEmail send verification url to user via email
func RequestVerifyEmail(
	getUserByID dependencyGetUserByID,
	upsertSession dependencyUpsertSession,
	sendEmail dependencySendEmail,
	tmpl *template.Template,
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
		token := uuid.NewV4().String()
		if err := upsertSession(models.Session{
			UserID: authToken.UserID,
			Token:  token,
			Type:   models.SessionTypeVerifyEmail,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// send email
		w := bytes.NewBufferString("")
		tmpl.Execute(w, map[string]interface{}{
			"email": url.QueryEscape(user.Email),
			"token": url.QueryEscape(token),
		})
		if err := sendEmail([]string{user.Email}, "SoleBTC --- Verify your email", w.String()); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
