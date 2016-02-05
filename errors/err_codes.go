package errors

// error codes
const (
	ErrCodeNil     Code = -1
	ErrCodeUnknown Code = 99999

	// validation error
	_ Code = iota + 4000
	ErrCodeInvalidBitcoinAddress

	// duplicate error
	_ Code = iota + 5000
	ErrCodeDuplicateEmail
	ErrCodeDuplicateBitcoinAddress
)
