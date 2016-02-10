package utils

import "fmt"

func ExampleBitcoinPrice() {
	_, err := bitcoinPriceWithByteFromBlockchain([]byte(`invalid-json-data`))
	fmt.Println(err != nil)

	price, _ := bitcoinPriceWithByteFromBlockchain([]byte(`{"USD": {"last": 5.5}}`))
	fmt.Println(price)

	// not able to test external service, just cover it
	BitcoinPrice()

	// Output:
	// true
	// 55000
}
