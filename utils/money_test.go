package utils

import "fmt"

func ExampleHumanReadableUSD() {
	fmt.Println(HumanReadableUSD(10))

	// Output:
	// 0.001
}

func ExampleMachineReadableUSD() {
	fmt.Println(MachineReadableUSD(0.01))

	// Output:
	// 100
}

func ExampleHumanReadableBTC() {
	fmt.Println(HumanReadableBTC(100000))

	// Output:
	// 0.001
}

func ExampleMachineReadableBTC() {
	fmt.Println(MachineReadableBTC(0.01))

	// Output:
	// 1000000
}
