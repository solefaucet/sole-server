package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleErrors() {
	errorString := "error"

	e := New(ErrCodeInvalidBitcoinAddress)
	e.ErrString = errorString
	e.ErrStringForLogging = "for logging purpose"
	fmt.Println(e, e.ErrCode)

	raw, _ := json.Marshal(e)
	fmt.Println(string(raw))

	// Output:
	// error 4001
	// {"error":"error"}
}
