package utils

import "fmt"

func ExampleToFixed() {
	fmt.Println(ToFixed(1.2345678, 0))
	fmt.Println(ToFixed(1.2345678, 1))
	fmt.Println(ToFixed(1.2345678, 2))
	fmt.Println(ToFixed(1.2345678, 3))

	// Output:
	// 1
	// 1.2
	// 1.23
	// 1.235
}
