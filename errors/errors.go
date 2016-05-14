package errors

import "errors"

// errors
var (
	ErrUnknown             = errors.New("unknown")
	ErrNotFound            = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrDuplicatedEmail     = errors.New("duplicated email")
	ErrDuplicatedAddress   = errors.New("duplicated address")
	ErrDuplicatedAuthToken = errors.New("duplicated auth token")
	ErrInvalidAddress      = errors.New("invalid address")
	ErrInvalidCaptcha      = errors.New("invalid captcha")
)
