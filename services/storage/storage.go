package storage

import (
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// Storage defines interface that one should implement
type Storage interface {
	// User
	CreateUser(models.User) *errors.Error

	// AuthToken
	GetAuthToken(authToken string) (models.AuthToken, *errors.Error)
	CreateAuthToken(models.AuthToken) *errors.Error
	DeleteAuthToken(authToken string) *errors.Error
}
