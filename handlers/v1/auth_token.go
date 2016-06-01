package v1

import (
	"net/http"

	"github.com/solefaucet/solebtc/errors"
	"github.com/solefaucet/solebtc/models"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type loginPayload struct {
	Email string `json:"email" binding:"required,email"`
}

// Login logs a existing user in, response with auth token
func Login(
	getUserByEmail dependencyGetUserByEmail,
	createAuthToken dependencyCreateAuthToken,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := loginPayload{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}

		user, err := getUserByEmail(payload.Email)
		if err != nil {
			switch err {
			case errors.ErrNotFound:
				c.AbortWithStatus(http.StatusNotFound)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		if user.Status == models.UserStatusBanned {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// create auth token with uuid v4
		authToken := models.AuthToken{
			UserID:    user.ID,
			AuthToken: uuid.NewV4().String(),
		}
		if err := createAuthToken(authToken); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, authToken)
	}
}

// Logout deletes corresponding auth token
func Logout(deleteAuthToken dependencyDeleteAuthToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)
		if err := deleteAuthToken(authToken.AuthToken); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
