package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleErrors() {
	errorString := "error"

	e := Error{
		ErrCode:             ErrCodeInvalidBitcoinAddress,
		ErrString:           errorString,
		ErrStringForLogging: "for logging purpose",
	}
	fmt.Println(e, e.ErrCode, e.NotNil())

	raw, _ := json.Marshal(e)
	fmt.Println(string(raw))

	// Output:
	// error 4003 true
	// {"error":"error"}
}
