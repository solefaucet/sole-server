package v1

import (
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// dependencies
type (
	// user
	dependencyGetUserByID    func(int64) (models.User, *errors.Error)
	dependencyGetUserByEmail func(string) (models.User, *errors.Error)
	dependencyCreateUser     func(models.User) *errors.Error
	dependencyUpdateUser     func(models.User) *errors.Error

	// auth token
	dependencyCreateAuthToken func(models.AuthToken) *errors.Error
	dependencyDeleteAuthToken func(string) *errors.Error

	// session
	dependencyUpsertSession     func(models.Session) *errors.Error
	dependencyGetSessionByToken func(string) (models.Session, *errors.Error)

	// email
	dependencySendEmail func(recipients []string, subject string, html string) *errors.Error
)
