package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/sole-server/errors"
	"github.com/solefaucet/sole-server/models"
)

type signupPayload struct {
	Email     string `json:"email" binding:"required,email"`
	Address   string `json:"address" binding:"required"`
	RefererID int64  `json:"referer_id,omitempty" binding:"-"`
}

func userWithSignupPayload(p signupPayload) models.User {
	return models.User{
		Email:     p.Email,
		Address:   strings.TrimSpace(p.Address),
		RefererID: p.RefererID,
	}
}

// Signup creates a new user with unique email, address
func Signup(
	validateAddress dependencyValidateAddress,
	createUser dependencyCreateUser,
	getUserByID dependencyGetUserByID,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := signupPayload{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}
		valid, _ := validateAddress(payload.Address)
		if !valid {
			c.AbortWithError(http.StatusBadRequest, errors.ErrInvalidAddress)
			return
		}

		user := userWithSignupPayload(payload)
		// assign referer_id to user
		referer, _ := getUserByID(payload.RefererID)
		user.RefererID = referer.ID

		if err := createUser(user); err != nil {
			switch err {
			case errors.ErrDuplicatedEmail:
				c.AbortWithError(http.StatusConflict, err)
			case errors.ErrDuplicatedAddress:
				c.AbortWithError(http.StatusConflict, err)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		logrus.WithFields(logrus.Fields{
			"event":   models.EventUserSignup,
			"email":   payload.Email,
			"address": payload.Address,
		}).Info("succeed to signup user")

		c.JSON(http.StatusOK, user)
	}
}

// VerifyEmail updates user's status to verified if current status is unverified
func VerifyEmail(
	getSessionByToken dependencyGetSessionByToken,
	getUserByID dependencyGetUserByID,
	updateUserStatus dependencyUpdateUserStatus,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")

		// check session lifetime
		session, err := getSessionByToken(token)
		if err != nil && err != errors.ErrNotFound {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if session.UpdatedAt.Add(3 * time.Hour).Before(time.Now()) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// get user
		user, err := getUserByID(session.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// check user status
		if user.Status == models.UserStatusBanned {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// update user
		if err := updateUserStatus(user.ID, models.UserStatusVerified); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

// UserInfo returns user's info as response
func UserInfo(getUserByID dependencyGetUserByID) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		user, err := getUserByID(authToken.UserID)
		if err != nil {
			// user is already authorized
			// if get user error
			// it must be internal server error
			// do not need to check existence of user
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// RefereeList returns user's referee list as response
func RefereeList(
	getReferees dependencyGetReferees,
	getNumberOfReferees dependencyGetNumberOfReferees,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse pagination args
		limit, offset, err := parsePagination(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		referees, err := getReferees(authToken.UserID, limit, offset)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		count, err := getNumberOfReferees(authToken.UserID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, paginationResult(referees, count))
	}
}
