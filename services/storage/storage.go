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
	CreateAuthToken(models.AuthToken) *errors.Error
}
