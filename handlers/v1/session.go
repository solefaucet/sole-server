package v1

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/solefaucet/sole-server/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// RequestVerifyEmail send verification url to user via email
func RequestVerifyEmail(
	getUserByID dependencyGetUserByID,
	upsertSession dependencyUpsertSession,
	sendEmail dependencySendEmail,
	tmpl *template.Template,
	appname string,
	appurl string,
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
			"appname": appname,
			"url":     appurl,
			"email":   url.QueryEscape(user.Email),
			"id":      user.ID,
			"token":   url.QueryEscape(token),
		})
		if err := sendEmail([]string{user.Email}, fmt.Sprintf("%s --- Verify your email", appname), w.String()); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
