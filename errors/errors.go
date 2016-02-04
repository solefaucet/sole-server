package errors

import "encoding/json"

var _ error = Error{}
var _ json.Marshaler = Error{}

// Error is custom error
type Error struct {
	ErrCode             code
	ErrString           string
	ErrStringForLogging string
}

type code int

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
