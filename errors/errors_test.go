package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleError() {
	e := New(ErrCodeInvalidBitcoinAddress)
	e.ErrStringForLogging = "for logging purpose"
	fmt.Println(e)

	raw, _ := json.Marshal(e)
	fmt.Println(string(raw))

	// Output:
	// 4001
	// {"code":4001}
}
