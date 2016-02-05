package storage

import (
	"github.com/freeusd/solebtc/errors"
	"github.com/freeusd/solebtc/models"
)

// Storage defines interface that one should implement
type Storage interface {
	CreateUser(models.User) *errors.Error
}
