package errors

import (
	"encoding/json"
	"fmt"
)

func ExampleError() {
	e := New(ErrCodeDuplicateEmail)
	e.ErrStringForLogging = "for logging purpose"
	fmt.Println(e)

	raw, _ := json.Marshal(e)
	fmt.Println(string(raw))

	// Output:
	// 5001
	// {"code":5001}
}
