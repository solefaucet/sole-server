package errors

import "encoding/json"

var (
	_ error          = New(0)
	_ json.Marshaler = New(0)
)

// Error is custom error
type Error struct {
	ErrCode             Code
	ErrString           string
	ErrStringForLogging string
}

// Code is custom error code
type Code int

// New creates new *Error with error code
func New(c Code) *Error {
	return &Error{
		ErrCode: c,
	}
}

// err implements error interface
func (e *Error) Error() string {
	return e.ErrString
}

// MarshalJSON implements json.Marshaler interface
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": e.Error(),
	})
}
