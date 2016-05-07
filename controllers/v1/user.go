package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
	"github.com/freeusd/solebtc/utils"
	"github.com/gin-gonic/gin"
)

// can be mocked out to test
var validateAddress = utils.ValidateAddress

type signupPayload struct {
	Email     string `json:"email" binding:"required,email"`
	Address   string `json:"address" binding:"required"`
	RefererID int64  `json:"referer_id,omitempty" binding:"-"`
}

func userWithSignupPayload(p signupPayload) models.User {
	return models.User{
		Email:   p.Email,
		Address: p.Address,
	}
}

// Signup creates a new user with unique email, address
func Signup(createUser dependencyCreateUser, getUserByID dependencyGetUserByID) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := signupPayload{}
		if err := c.BindJSON(&payload); err != nil {
			return
		}
		valid, err := validateAddress(payload.Address)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !valid {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s is invalid address", payload.Address))
			return
		}

		user := userWithSignupPayload(payload)
		// assign referer_id to user
		referer, _ := getUserByID(payload.RefererID)
		user.RefererID = referer.ID

		if err := createUser(user); err != nil {
			switch err.ErrCode {
			case errors.ErrCodeDuplicateEmail:
				c.AbortWithError(http.StatusConflict, err)
			case errors.ErrCodeDuplicateAddress:
				c.AbortWithError(http.StatusConflict, err)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

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
		if err != nil {
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
			// although error code can be ErrCodeNotFound
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// RefereeList returns user's referee list as response
func RefereeList(
	getRefereesSinceID dependencyGetRefereesSince,
	getRefereesUntilID dependencyGetRefereesUntil,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.MustGet("auth_token").(models.AuthToken)

		// parse pagination args
		isSince, separator, limit, err := parsePagination(c)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// get result according to args
		result := []models.User{}
		var syserr *errors.Error
		if isSince {
			result, syserr = getRefereesSinceID(authToken.UserID, separator, limit)
		} else {
			result, syserr = getRefereesUntilID(authToken.UserID, separator, limit)
		}

		// response with result or error
		if syserr != nil {
			c.AbortWithError(http.StatusInternalServerError, syserr)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
