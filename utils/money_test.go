package utils

import "fmt"

func ExampleInitializePriceConverter() {
	for _, typ := range []string{"usd", "btc", "unknown"} {
		fmt.Println(InitializePriceConverter(typ))
	}

	// Output:
	// <nil>
	// <nil>
	// converter unknown not exist
}

func ExampleToMachine() {
	InitializePriceConverter("usd")
	fmt.Println(ToMachine(0.01))

	InitializePriceConverter("btc")
	fmt.Println(machineReadableBTC(0.01))

	// Output:
	// 100
	// 1000000
}

func ExampleToHuman() {
	InitializePriceConverter("usd")
	fmt.Println(humanReadableUSD(10))

	InitializePriceConverter("btc")
	fmt.Println(ToHuman(100000))

	// Output:
	// 0.001
	// 0.001
}
