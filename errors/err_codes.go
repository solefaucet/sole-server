package errors

// error codes
const (
	_ Code = iota
	ErrCodeUnknown
	ErrCodeNotFound
	ErrCodeInsufficientBalance
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

// external service error
const (
	_ Code = iota + 3000
	ErrCodeMandrill
)
