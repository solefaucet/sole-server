package utils

import "fmt"

func ExampleBitcoinPrice() {
	_, err := bitcoinPriceWithByteFromBlockchain([]byte(`invalid-json-data`))
	fmt.Println(err != nil)

	_, err = bitcoinPriceWithByteFromBlockchain([]byte(`{}`))
	fmt.Println(err)

	price, _ := bitcoinPriceWithByteFromBlockchain([]byte(`{"USD": {"last": 5.5}}`))
	fmt.Println(price)

	// not able to test external service, just cover it
	BitcoinPrice()

	// Output:
	// true
	// bitcoin price cannot be 0
	// 55000
}
