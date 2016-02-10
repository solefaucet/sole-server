package v1

import (
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

func mockGetUserByEmail(user models.User, err *errors.Error) dependencyGetUserByEmail {
	return func(string) (models.User, *errors.Error) {
		return user, err
	}
}

func mockGetUserByID(user models.User, err *errors.Error) dependencyGetUserByID {
	return func(int64) (models.User, *errors.Error) {
		return user, err
	}
}

func mockCreateUser(err *errors.Error) dependencyCreateUser {
	return func(models.User) *errors.Error {
		return err
	}
}

func mockUpdateUser(err *errors.Error) dependencyUpdateUser {
	return func(models.User) *errors.Error {
		return err
	}
}

func mockCreateAuthToken(err *errors.Error) dependencyCreateAuthToken {
	return func(models.AuthToken) *errors.Error {
		return err
	}
}

func mockDeleteAuthToken(err *errors.Error) dependencyDeleteAuthToken {
	return func(string) *errors.Error {
		return err
	}
}

func mockGetSessionByToken(sess models.Session, err *errors.Error) dependencyGetSessionByToken {
	return func(string) (models.Session, *errors.Error) {
		return sess, err
	}
}

func mockUpsertSession(err *errors.Error) dependencyUpsertSession {
	return func(models.Session) *errors.Error {
		return err
	}
}

func mockSendEmail(err *errors.Error) dependencySendEmail {
	return func([]string, string, string) *errors.Error {
		return err
	}
}
