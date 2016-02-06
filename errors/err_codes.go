package errors

// error codes
const (
	_ Code = iota
	ErrCodeUnknown
)

// duplicate error
const (
	_ Code = iota + 5000
	ErrCodeDuplicateEmail
	ErrCodeDuplicateBitcoinAddress
	ErrCodeDuplicateAuthToken
)

// validate error
const (
	_ Code = iota + 4000
	ErrCodeInvalidBitcoinAddress
)
