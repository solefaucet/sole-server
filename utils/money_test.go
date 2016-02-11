package utils

import "fmt"

func ExampleMoney() {
	fmt.Println(MachineReadableUSD(0.01))
	fmt.Println(HumanReadableUSD(10))
	fmt.Println(MachineReadableBTC(0.01))
	fmt.Println(HumanReadableBTC(100000))

	// Output:
	// 100
	// 0.001
	// 1000000
	// 0.001
}
