package errors

import "encoding/json"

var (
	_ error          = Error{}
	_ json.Marshaler = Error{}
	// Nil error
	Nil = Error{ErrCode: ErrCodeNil}
)

// Error is custom error
type Error struct {
	ErrCode             Code
	ErrString           string
	ErrStringForLogging string
}

// Code is custom error code
type Code int

// NotNil checks if error is not nil
func (e Error) NotNil() bool {
	return e.ErrCode != ErrCodeNil
}

// err implements error interface
func (e Error) Error() string {
	return e.ErrString
}

// MarshalJSON implements json.Marshaler interface
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": e.Error(),
	})
}
