package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleErrors() {
	errorString := "error"

	e := Error{
		ErrCode:   ErrCodeInvalidBitcoinAddress,
		ErrString: errorString,
	}
	fmt.Println(e, e.ErrCode)

	raw, _ := json.Marshal(e)
	fmt.Println(string(raw))

	// Output:
	// error 40001
	// {"error":"error"}
}
